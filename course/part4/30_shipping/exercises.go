// Lesson 30 exercises: shipping. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
// Run the whole app:              DATABASE_PATH=/tmp/x.db go run ./cmd/app
package lesson30

import (
	"context"
	"net"
	"net/http"
)

// Config holds the app's settings.
type Config struct {
	Port         int    // from PORT, default 8080, must be a number between 1 and 65535
	LogLevel     string // from LOG_LEVEL, default "info", must be one of: debug, info, warn, error
	DatabasePath string // from DATABASE_PATH, REQUIRED
}

// Ex1: read the settings using getenv (os.Getenv in main, a fake in tests). (README section 2.)
//   - PORT: empty → 8080; not a number, or outside 1..65535 → error "invalid PORT: <value>"
//   - LOG_LEVEL: empty → "info"; anything other than debug/info/warn/error → error "invalid LOG_LEVEL: <value>"
//   - DATABASE_PATH: empty → error "DATABASE_PATH is required"
//
// On any error, return Config{} and the error.
func LoadConfig(getenv func(string) string) (Config, error) {
	// TODO
	return Config{}, nil
}

// Ex2: embed the file static/index.html into this variable. (README section 3.)
// Add the //go:embed line DIRECTLY above the var (and import _ "embed" at the top of the file).
var indexHTML string

// Index serves the embedded page. Already done: it works once indexHTML is filled in.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

// Version is set at build time with:  go build -ldflags "-X gomastery/course/part4/30_shipping.Version=1.2.3"
// (cmd/app has its own main.Version that it copies here.)
var Version = "dev"

// Ex3: respond 200 with the JSON {"status":"ok","version":"<Version>"} and Content-Type application/json.
// Read the Version variable when the request arrives (not at startup).
func Health(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// Ex4: serve h on the listener ln until ctx is cancelled, then shut down GRACEFULLY:
// requests that are already running get up to 5 seconds to finish. (README section 5.)
//
//  1. srv := &http.Server{Handler: h}
//  2. start a goroutine that waits for <-ctx.Done(), then calls srv.Shutdown with a NEW
//     context.WithTimeout(context.Background(), 5*time.Second)
//  3. err := srv.Serve(ln)    (Serve is like ListenAndServe, on a listener you already have)
//  4. if errors.Is(err, http.ErrServerClosed), that's a normal shutdown: return nil
//     (but first wait for Shutdown to finish, so in-flight requests complete: use a done channel)
//  5. any other error: return it
func Serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	// TODO
	return nil
}
