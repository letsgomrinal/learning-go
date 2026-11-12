package lesson26

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// Always run these with:  go test -race -v

// do sends a request to the server (no real network) and returns the recorded response.
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
		t.Fatalf("the response body isn't valid JSON: %q (%v)", rec.Body.String(), err)
	}
	return v
}

func wantStatus(t *testing.T, rec *httptest.ResponseRecorder, want int, what string) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("%s: status %d, want %d (body %q)", what, rec.Code, want, rec.Body.String())
	}
}

func TestEx1Helpers(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusTeapot, map[string]int{"n": 1})
	if rec.Code != 418 || rec.Header().Get("Content-Type") != "application/json" || strings.TrimSpace(rec.Body.String()) != `{"n":1}` {
		t.Errorf("Ex1: writeJSON: status %d, Content-Type %q, body %q; want 418, application/json, {\"n\":1}",
			rec.Code, rec.Header().Get("Content-Type"), rec.Body.String())
	}
	rec = httptest.NewRecorder()
	writeError(rec, http.StatusBadRequest, "nope")
	if rec.Code != 400 || strings.TrimSpace(rec.Body.String()) != `{"error":"nope"}` {
		t.Errorf("Ex1: writeError: status %d, body %q; want 400, {\"error\":\"nope\"}", rec.Code, rec.Body.String())
	}
}

func TestEx2Health(t *testing.T) {
	rec := do(t, NewServer(), "GET", "/health", "")
	wantStatus(t, rec, 200, "Ex2: GET /health")
	if got := decode[map[string]string](t, rec); got["status"] != "ok" {
		t.Errorf("Ex2: body = %v, want {\"status\":\"ok\"}", got)
	}
}

func TestEx3CreateNote(t *testing.T) {
	s := NewServer()
	rec := do(t, s, "POST", "/notes", `{"text": "  buy milk "}`)
	wantStatus(t, rec, 201, "Ex3: POST /notes")
	if got := decode[Note](t, rec); got != (Note{ID: 1, Text: "buy milk"}) {
		t.Errorf("Ex3: created %+v, want {ID:1 Text:buy milk}", got)
	}
	rec = do(t, s, "POST", "/notes", `{"text": "call mom"}`)
	if got := decode[Note](t, rec); got.ID != 2 {
		t.Errorf("Ex3: second note has id %d, want 2", got.ID)
	}
	for _, body := range []string{`not json`, `{"text": "   "}`, `{}`} {
		rec := do(t, s, "POST", "/notes", body)
		wantStatus(t, rec, 400, fmt.Sprintf("Ex3: POST /notes with %s", body))
		if decode[map[string]string](t, rec)["error"] == "" {
			t.Errorf("Ex3: error responses should look like {\"error\": \"...\"}")
		}
	}
}

func TestEx4GetNote(t *testing.T) {
	s := NewServer()
	do(t, s, "POST", "/notes", `{"text": "first"}`)
	rec := do(t, s, "GET", "/notes/1", "")
	wantStatus(t, rec, 200, "Ex4: GET /notes/1")
	if got := decode[Note](t, rec); got != (Note{ID: 1, Text: "first"}) {
		t.Errorf("Ex4: got %+v", got)
	}
	wantStatus(t, do(t, s, "GET", "/notes/99", ""), 404, "Ex4: GET /notes/99")
	wantStatus(t, do(t, s, "GET", "/notes/abc", ""), 400, "Ex4: GET /notes/abc")
}

func TestEx5ListNotes(t *testing.T) {
	s := NewServer()
	rec := do(t, s, "GET", "/notes", "")
	wantStatus(t, rec, 200, "Ex5: GET /notes (empty)")
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("Ex5: empty list body = %q, want []", rec.Body.String())
	}
	for _, text := range []string{"a", "b", "c"} {
		do(t, s, "POST", "/notes", fmt.Sprintf(`{"text": %q}`, text))
	}
	notes := decode[[]Note](t, do(t, s, "GET", "/notes", ""))
	if len(notes) != 3 || notes[0].ID != 1 || notes[1].ID != 2 || notes[2].Text != "c" {
		t.Errorf("Ex5: GET /notes = %+v, want 3 notes sorted by id", notes)
	}
}

func TestEx6DeleteNote(t *testing.T) {
	s := NewServer()
	do(t, s, "POST", "/notes", `{"text": "temp"}`)
	rec := do(t, s, "DELETE", "/notes/1", "")
	wantStatus(t, rec, 204, "Ex6: DELETE /notes/1")
	if rec.Body.Len() != 0 {
		t.Errorf("Ex6: a 204 response must have no body, got %q", rec.Body.String())
	}
	wantStatus(t, do(t, s, "GET", "/notes/1", ""), 404, "Ex6: GET after DELETE")
	wantStatus(t, do(t, s, "DELETE", "/notes/1", ""), 404, "Ex6: DELETE twice")
	wantStatus(t, do(t, s, "DELETE", "/notes/x", ""), 400, "Ex6: DELETE /notes/x")
}

func TestEx7Search(t *testing.T) {
	s := NewServer()
	for _, text := range []string{"Buy MILK", "walk dog", "milkshake recipe"} {
		do(t, s, "POST", "/notes", fmt.Sprintf(`{"text": %q}`, text))
	}
	notes := decode[[]Note](t, do(t, s, "GET", "/notes?q=milk", ""))
	if len(notes) != 2 || notes[0].ID != 1 || notes[1].ID != 3 {
		t.Errorf("Ex7: GET /notes?q=milk = %+v, want notes 1 and 3", notes)
	}
	if notes := decode[[]Note](t, do(t, s, "GET", "/notes?q=zzz", "")); len(notes) != 0 {
		t.Errorf("Ex7: no matches should give [], got %+v", notes)
	}
}

func TestConcurrentCreates(t *testing.T) {
	// Requests run concurrently in a real server. This catches a missing lock (with -race).
	s := NewServer()
	if do(t, s, "POST", "/notes", `{"text":"probe"}`).Code != 201 {
		t.Skip("finish Ex3 first")
	}
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			do(t, s, "POST", "/notes", fmt.Sprintf(`{"text":"n%d"}`, i))
		}()
	}
	wg.Wait()
	if rec := do(t, s, "GET", "/notes", ""); rec.Code == 200 {
		if notes := decode[[]Note](t, rec); len(notes) != 51 {
			t.Errorf("after 51 creates, GET /notes has %d notes (lost writes?)", len(notes))
		}
	}
}
