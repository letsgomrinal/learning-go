package lesson27

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

func TestEx1NewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("Ex1: NewClient returned nil")
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("Ex1: Timeout = %v, want 5s", c.Timeout)
	}
}

// fakeAPI is a small fake user service for the client tests.
func fakeAPI(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "hello there") })
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done():
		}
	})
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") != "7" {
			http.Error(w, "no such user", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id": 7, "name": "Ana", "email": "ana@example.com"}`)
	})
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "send JSON", http.StatusUnsupportedMediaType)
			return
		}
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil || u.Name == "" {
			http.Error(w, "bad user", http.StatusBadRequest)
			return
		}
		u.ID = 42
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestEx2GetText(t *testing.T) {
	srv := fakeAPI(t)
	ctx := context.Background()
	if got, err := GetText(ctx, srv.Client(), srv.URL+"/hello"); got != "hello there" || err != nil {
		t.Errorf("Ex2: GetText(/hello) = %q, %v; want \"hello there\", nil", got, err)
	}
	if _, err := GetText(ctx, srv.Client(), srv.URL+"/nope"); err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("Ex2: GetText(/nope) error = %v, want \"unexpected status 404\"", err)
	}
	ctx2, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	if _, err := GetText(ctx2, srv.Client(), srv.URL+"/slow"); err == nil {
		t.Error("Ex2: GetText with a cancelled context should return an error")
	}
	if time.Since(start) > time.Second {
		t.Error("Ex2: GetText ignored the context's deadline. Use http.NewRequestWithContext")
	}
}

func TestEx3GetUser(t *testing.T) {
	srv := fakeAPI(t)
	u, err := GetUser(context.Background(), srv.Client(), srv.URL, 7)
	if err != nil || u != (User{ID: 7, Name: "Ana", Email: "ana@example.com"}) {
		t.Errorf("Ex3: GetUser(7) = %+v, %v", u, err)
	}
	if _, err := GetUser(context.Background(), srv.Client(), srv.URL, 8); err == nil {
		t.Error("Ex3: GetUser(8) should return an error (the server says 404)")
	}
}

func TestEx4CreateUser(t *testing.T) {
	srv := fakeAPI(t)
	u, err := CreateUser(context.Background(), srv.Client(), srv.URL, User{Name: "Bo", Email: "bo@example.com"})
	if err != nil || u != (User{ID: 42, Name: "Bo", Email: "bo@example.com"}) {
		t.Errorf("Ex4: CreateUser = %+v, %v; want ID 42 assigned by the server", u, err)
	}
	if _, err := CreateUser(context.Background(), srv.Client(), srv.URL, User{}); err == nil {
		t.Error("Ex4: the server rejects a user with no name (400): CreateUser should return an error")
	}
}

var ok200 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "ok") })

func serve(h http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestEx5Logging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	rec := serve(Logging(logger)(ok200), "GET", "/notes", nil)
	if rec.Body.String() != "ok" {
		t.Errorf("Ex5: the wrapped handler didn't run (body %q). Call next.ServeHTTP(w, r)", rec.Body.String())
	}
	logged := buf.String()
	for _, want := range []string{"msg=request", "method=GET", "path=/notes"} {
		if !strings.Contains(logged, want) {
			t.Errorf("Ex5: log output %q doesn't contain %q", logged, want)
		}
	}
}

func TestEx6Recover(t *testing.T) {
	boom := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("something broke") })
	var rec *httptest.ResponseRecorder
	func() {
		defer func() {
			if v := recover(); v != nil {
				t.Fatalf("Ex6: the panic escaped the middleware: %v", v)
			}
		}()
		rec = serve(Recover(boom), "GET", "/", nil)
	}()
	if rec.Code != 500 || strings.TrimSpace(rec.Body.String()) != `{"error":"internal error"}` {
		t.Errorf("Ex6: status %d, body %q; want 500, {\"error\":\"internal error\"}", rec.Code, rec.Body.String())
	}
	if rec := serve(Recover(ok200), "GET", "/", nil); rec.Body.String() != "ok" {
		t.Errorf("Ex6: without a panic the handler should run normally, got %q", rec.Body.String())
	}
}

func TestEx7RequireAPIKey(t *testing.T) {
	h := RequireAPIKey("s3cret")(ok200)
	if rec := serve(h, "GET", "/", map[string]string{"X-API-Key": "s3cret"}); rec.Code != 200 || rec.Body.String() != "ok" {
		t.Errorf("Ex7: correct key: status %d, body %q; want 200, ok", rec.Code, rec.Body.String())
	}
	for _, headers := range []map[string]string{nil, {"X-API-Key": "wrong"}} {
		rec := serve(h, "GET", "/", headers)
		if rec.Code != 401 || strings.TrimSpace(rec.Body.String()) != `{"error":"unauthorized"}` {
			t.Errorf("Ex7: key %q: status %d, body %q; want 401, {\"error\":\"unauthorized\"}", headers["X-API-Key"], rec.Code, rec.Body.String())
		}
	}
}

func TestEx8Chain(t *testing.T) {
	var order []string
	mw := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { order = append(order, "handler") })
	serve(Chain(final, mw("a"), mw("b"), mw("c")), "GET", "/", nil)
	if got := strings.Join(order, " "); got != "a b c handler" {
		t.Errorf("Ex8: order = %q, want \"a b c handler\"", got)
	}
}
