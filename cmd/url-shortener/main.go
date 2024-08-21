package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	config "restapi/iternal"
	"restapi/iternal/http-server/handlers/url/delete"
	"restapi/iternal/http-server/handlers/url/redirect"
	"restapi/iternal/http-server/handlers/url/save"
	mwLogger "restapi/iternal/http-server/middleware/logger"
	"restapi/iternal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting server...", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	db, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		slog.Error("failed to open database: %w", err)
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			slog.Info("failed to close database: %w", err)
		}
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.BasicAuth("basic-auth", map[string]string{
			cfg.HttpServer.User: cfg.HttpServer.Password,
		}))
		r.Post("/create", save.New(log, db))
		r.Delete("/delete/{alias}", delete.New(log, db))
	})

	router.Get("/{alias}", redirect.New(log, db))

	log.Info("Server started",
		slog.String("address", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		slog.Error("failed to start server: %w", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))
	}
	return log
}
