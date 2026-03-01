package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/handlers"
	"github.com/CameronBrooks11/cameronbrooks-site/internal/middleware"
	"github.com/CameronBrooks11/cameronbrooks-site/static"
)

// Version and BuildTime are injected at compile time via -ldflags.
// During go run (development) they hold their zero values.
//
//	make build
//	-> -ldflags="-X main.Version=<git-sha> -X main.BuildTime=<iso8601>"
var (
	Version   string
	BuildTime string
)

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("SITE_ENV") == "dev" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	if err := handlers.InitTemplates(); err != nil {
		slog.Error("failed to initialize templates", "err", err)
		os.Exit(1)
	}

	h := handlers.New()
	h.AppVersion = Version
	h.AppBuildTime = BuildTime

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /projects", h.Projects)
	mux.HandleFunc("GET /projects/{slug}", h.Project)
	mux.HandleFunc("GET /writing", h.Writing)
	mux.HandleFunc("GET /writing/{slug}", h.Post)
	mux.HandleFunc("GET /about", h.About)
	mux.HandleFunc("GET /contact", h.Contact)

	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /version", h.Version)

	mux.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))),
	)

	addr := os.Getenv("SITE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      middleware.Chain(mux, middleware.RequestID, middleware.Logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		slog.Info("starting server", "addr", srv.Addr, "version", Version)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("shutting down", "reason", ctx.Err())
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
	slog.Info("shutdown complete")
}
