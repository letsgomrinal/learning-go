package lesson30

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

// env builds a fake getenv from a map.
func env(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}

func TestEx1LoadConfig(t *testing.T) {
	cfg, err := LoadConfig(env(map[string]string{"DATABASE_PATH": "/data/app.db"}))
	if err != nil || cfg != (Config{Port: 8080, LogLevel: "info", DatabasePath: "/data/app.db"}) {
		t.Errorf("Ex1: defaults: got %+v, %v; want {8080 info /data/app.db}", cfg, err)
	}
	cfg, err = LoadConfig(env(map[string]string{"PORT": "9000", "LOG_LEVEL": "debug", "DATABASE_PATH": "x.db"}))
	if err != nil || cfg != (Config{Port: 9000, LogLevel: "debug", DatabasePath: "x.db"}) {
		t.Errorf("Ex1: all set: got %+v, %v", cfg, err)
	}
	bad := []struct {
		vars map[string]string
		want string
	}{
		{map[string]string{}, "DATABASE_PATH is required"},
		{map[string]string{"DATABASE_PATH": "x", "PORT": "eighty"}, "invalid PORT: eighty"},
		{map[string]string{"DATABASE_PATH": "x", "PORT": "70000"}, "invalid PORT: 70000"},
		{map[string]string{"DATABASE_PATH": "x", "LOG_LEVEL": "loud"}, "invalid LOG_LEVEL: loud"},
	}
	for _, b := range bad {
		cfg, err := LoadConfig(env(b.vars))
		if err == nil || err.Error() != b.want || cfg != (Config{}) {
			t.Errorf("Ex1: LoadConfig(%v) = %+v, %v; want Config{} and %q", b.vars, cfg, err, b.want)
		}
	}
}

func TestEx2Embed(t *testing.T) {
	if !strings.Contains(indexHTML, "<h1>Hello from an embedded file</h1>") {
		t.Fatalf("Ex2: indexHTML is %d characters and doesn't contain the page. Add //go:embed static/index.html above the var", len(indexHTML))
	}
	rec := httptest.NewRecorder()
	Index(rec, httptest.NewRequest("GET", "/", nil))
	if !strings.Contains(rec.Body.String(), "<h1>") {
		t.Errorf("Ex2: Index served %q", rec.Body.String())
	}
}

func TestEx3Health(t *testing.T) {
	old := Version
	Version = "1.2.3"
	defer func() { Version = old }()

	rec := httptest.NewRecorder()
	Health(rec, httptest.NewRequest("GET", "/health", nil))
	var body map[string]string
	json.Unmarshal(rec.Body.Bytes(), &body)
	if rec.Code != 200 || rec.Header().Get("Content-Type") != "application/json" || body["status"] != "ok" || body["version"] != "1.2.3" {
		t.Errorf("Ex3: status %d, Content-Type %q, body %q; want 200, application/json, {\"status\":\"ok\",\"version\":\"1.2.3\"}",
			rec.Code, rec.Header().Get("Content-Type"), rec.Body.String())
	}
}

func TestEx4GracefulShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0") // port 0: pick any free port
	if err != nil {
		t.Fatal(err)
	}
	slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond) // a request that's still running during shutdown
		fmt.Fprint(w, "finished")
	})

	ctx, cancel := context.WithCancel(context.Background())
	serveErr := make(chan error, 1)
	go func() { serveErr <- Serve(ctx, ln, slow) }()

	// Start a slow request, then ask the server to stop while it's running.
	type result struct {
		body string
		err  error
	}
	got := make(chan result, 1)
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String() + "/")
		if err != nil {
			got <- result{err: err}
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		got <- result{body: string(b)}
	}()
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-serveErr:
		if err != nil {
			t.Errorf("Ex4: Serve returned %v after a clean shutdown, want nil (http.ErrServerClosed is normal)", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Ex4: Serve didn't return after ctx was cancelled")
	}
	select {
	case r := <-got:
		if r.err != nil || r.body != "finished" {
			t.Errorf("Ex4: the in-flight request got %q, %v; want \"finished\". Shut down gracefully with srv.Shutdown", r.body, r.err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Ex4: the in-flight request never completed")
	}
}
