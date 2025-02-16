package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"testFojune/internal/config"
	"testFojune/internal/db/initdb"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/handlers/change"
	deletewallet "testFojune/internal/http-server/handlers/delete"
	"testFojune/internal/http-server/handlers/get"
	"testFojune/internal/http-server/handlers/save"
)

func main() {
	cfg := config.NewConfig()

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, opts))

	log.Info("starting")
	log.Debug("enable Debug level")

	db, err := initdb.InitDB(config.DB{
		User:     cfg.User,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
	})
	if err != nil {
		log.Error("failed to init storage", errlog.Err(err))
		os.Exit(1)
	}

	log.Info("storage is initialized")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// save wallet
	router.Post("/wallet/save", save.New(log, db))
	// get wallet balance
	router.Get("/wallet/{uuid}", get.New(log, db))
	// change wallet balance
	router.Post("/wallet", change.New(log, db))
	// delete wallet
	router.Delete("/wallet/delete", deletewallet.New(log, db))

	log.Info("starting server")

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	// TODO: server
}
