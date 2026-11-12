package todoapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// do sends a request straight to the handler (no network) and returns the recorder.
func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("invalid JSON body %q: %v", rec.Body.String(), err)
	}
	return v
}

func TestCreateAndGet(t *testing.T) {
	s := NewServer()

	rec := do(t, s, "POST", "/todos", `{"title":"learn go"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want 201; body %q", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q", ct)
	}
	created := decode[Todo](t, rec)
	if created != (Todo{ID: 1, Title: "learn go"}) {
		t.Errorf("created = %+v", created)
	}

	rec = do(t, s, "GET", "/todos/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET status = %d", rec.Code)
	}
	if got := decode[Todo](t, rec); got != created {
		t.Errorf("GET = %+v, want %+v", got, created)
	}
}

func TestList(t *testing.T) {
	s := NewServer()

	rec := do(t, s, "GET", "/todos", "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("empty list: status %d body %q, want 200 []", rec.Code, rec.Body.String())
	}

	for _, title := range []string{"a", "b", "c"} {
		do(t, s, "POST", "/todos", fmt.Sprintf(`{"title":%q}`, title))
	}
	todos := decode[[]Todo](t, do(t, s, "GET", "/todos", ""))
	if len(todos) != 3 || todos[0].ID != 1 || todos[2].Title != "c" {
		t.Errorf("list = %+v", todos)
	}
}

func TestValidation(t *testing.T) {
	s := NewServer()
	cases := []struct {
		method, path, body string
		want               int
	}{
		{"POST", "/todos", `not json`, 400},
		{"POST", "/todos", `{"title":"   "}`, 400},
		{"GET", "/todos/abc", "", 400},
		{"GET", "/todos/99", "", 404},
		{"PATCH", "/todos/99", `{"done":true}`, 404},
		{"DELETE", "/todos/99", "", 404},
		{"PUT", "/todos", "", 405}, // ServeMux handles this for you
	}
	for _, c := range cases {
		rec := do(t, s, c.method, c.path, c.body)
		if rec.Code != c.want {
			t.Errorf("%s %s: status %d, want %d", c.method, c.path, rec.Code, c.want)
			continue
		}
		if c.want == 400 || c.want == 404 {
			body := decode[map[string]string](t, rec)
			if body["error"] == "" {
				t.Errorf("%s %s: want JSON {\"error\": ...}, got %q", c.method, c.path, rec.Body.String())
			}
		}
	}
}

func TestPatchAndDelete(t *testing.T) {
	s := NewServer()
	do(t, s, "POST", "/todos", `{"title":"x"}`)

	rec := do(t, s, "PATCH", "/todos/1", `{"done":true}`)
	if rec.Code != 200 {
		t.Fatalf("PATCH status = %d", rec.Code)
	}
	if got := decode[Todo](t, rec); !got.Done || got.Title != "x" {
		t.Errorf("PATCH = %+v", got)
	}

	if rec := do(t, s, "DELETE", "/todos/1", ""); rec.Code != http.StatusNoContent {
		t.Errorf("DELETE status = %d, want 204", rec.Code)
	}
	if rec := do(t, s, "GET", "/todos/1", ""); rec.Code != 404 {
		t.Errorf("GET after DELETE = %d, want 404", rec.Code)
	}
}

func TestConcurrentCreates(t *testing.T) {
	// A real server: requests go over the network and every request is its own goroutine.
	s := NewServer()
	srv := httptest.NewServer(s)
	defer srv.Close()

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(srv.URL+"/todos", "application/json",
				strings.NewReader(fmt.Sprintf(`{"title":"t%d"}`, i)))
			if err != nil {
				t.Error(err)
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()

	todos := decode[[]Todo](t, do(t, s, "GET", "/todos", ""))
	if len(todos) != 50 {
		t.Fatalf("got %d todos, want 50 (lost writes? run with -race)", len(todos))
	}
	seen := map[int]bool{}
	for _, td := range todos {
		if seen[td.ID] {
			t.Fatalf("duplicate id %d", td.ID)
		}
		seen[td.ID] = true
	}
}
