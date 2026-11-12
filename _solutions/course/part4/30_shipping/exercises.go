// Reference solution for lesson 30.
package lesson30

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Port         int
	LogLevel     string
	DatabasePath string
}

func LoadConfig(getenv func(string) string) (Config, error) {
	cfg := Config{Port: 8080, LogLevel: "info", DatabasePath: getenv("DATABASE_PATH")}

	if p := getenv("PORT"); p != "" {
		port, err := strconv.Atoi(p)
		if err != nil || port < 1 || port > 65535 {
			return Config{}, fmt.Errorf("invalid PORT: %s", p)
		}
		cfg.Port = port
	}

	if l := getenv("LOG_LEVEL"); l != "" {
		switch l {
		case "debug", "info", "warn", "error":
			cfg.LogLevel = l
		default:
			return Config{}, fmt.Errorf("invalid LOG_LEVEL: %s", l)
		}
	}

	if cfg.DatabasePath == "" {
		return Config{}, errors.New("DATABASE_PATH is required")
	}
	return cfg, nil
}

//go:embed static/index.html
var indexHTML string

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

var Version = "dev"

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": Version})
}

func Serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	srv := &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second}

	shutdownDone := make(chan error, 1)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownDone <- srv.Shutdown(shutdownCtx)
	}()

	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return <-shutdownDone // wait until in-flight requests have finished
	}
	return err
}
