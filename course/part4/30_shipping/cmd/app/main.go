// Command app wires lesson 30 together: config, embedded page, health check, graceful shutdown.
//
//	DATABASE_PATH=/tmp/x.db go run -ldflags "-X main.Version=0.1.0" ./cmd/app
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	lesson30 "gomastery/course/part4/30_shipping"
)

// Version is set at build time: -ldflags "-X main.Version=..."
var Version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := lesson30.LoadConfig(os.Getenv)
	if err != nil {
		return err // fail fast, with a clear message
	}
	lesson30.Version = Version

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", lesson30.Index)
	mux.HandleFunc("GET /health", lesson30.Health)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return err
	}
	logger.Info("listening", "port", cfg.Port, "version", Version, "log_level", cfg.LogLevel)
	if err := lesson30.Serve(ctx, ln, mux); err != nil {
		return err
	}
	logger.Info("shut down cleanly")
	return nil
}
