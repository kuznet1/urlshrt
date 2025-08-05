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
	svc := service.NewService(&repository.MemoryRepo{})
	h := handler.NewHandler(svc, cfg)
	mux := chi.NewRouter()
	mux.Post("/", middleware.Logging(h.Shorten))
	mux.Get("/{id}", middleware.Logging(h.Lengthen))

	fmt.Println("Shortener service is starting at", cfg.ListenAddr)
	err := http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		panic(err)
	}
}
