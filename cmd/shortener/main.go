package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/handler"
	"github.com/kuznet1/urlshrt/internal/middleware"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"net/http"
)

func main() {
	cfg := config.ParseArgs()
	repo, err := repository.NewMemoryRepo(cfg.FileStoragePath)
	if err != nil {
		panic(err)
	}
	svc := service.NewService(repo, cfg)
	h := handler.NewHandler(svc)
	mux := chi.NewRouter()
	mux.Post("/", middleware.Logging(middleware.Compression(h.Shorten)))
	mux.Get("/{id}", middleware.Logging(middleware.Compression(h.Lengthen)))
	mux.Post("/api/shorten", middleware.Logging(middleware.Compression(h.ShortenJSON)))

	fmt.Println("Shortener service is starting at", cfg.ListenAddr)
	err = http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		panic(err)
	}
}
