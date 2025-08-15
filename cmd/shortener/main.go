package main

import (
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
)

func main() {
	cfg, err := config.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewMemoryRepo(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(repo, cfg)
	h := handler.NewHandler(svc, logger)
	mux := chi.NewRouter()

	requestLogger := middleware.NewRequestLogger(logger)

	mux.Post("/", requestLogger.Wrap(middleware.Compression(h.Shorten)))
	mux.Get("/{id}", requestLogger.Wrap(middleware.Compression(h.Lengthen)))
	mux.Post("/api/shorten", requestLogger.Wrap(middleware.Compression(h.ShortenJSON)))

	fmt.Println("Shortener service is starting at", cfg.ListenAddr)
	err = http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
