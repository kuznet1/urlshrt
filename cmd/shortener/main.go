package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/handler"
	"github.com/kuznet1/urlshrt/internal/middleware"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"go.uber.org/zap"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewMemoryRepo(cfg.FileStoragePath, logger)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(repo, cfg)
	h := handler.NewHandler(svc, logger)
	requestLogger := middleware.NewRequestLogger(logger)
	mux := chi.NewRouter()
	mux.Use(requestLogger.Logging, middleware.Compression)
	mux.Post("/", h.Shorten)
	mux.Get("/{id}", h.Lengthen)
	mux.Post("/api/shorten", h.ShortenJSON)
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			logger.Error("db conn error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Shortener service is starting at", cfg.ListenAddr)
	err = http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
