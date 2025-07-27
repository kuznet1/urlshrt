package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/handler"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"net/http"
)

func main() {
	svc := service.NewService(&repository.MemoryRepo{})
	h := handler.NewHandler(svc)
	mux := chi.NewRouter()
	mux.Post("/", h.Shorten)
	mux.Get("/{id}", h.Lengthen)

	fmt.Println("Shortener service is starting at :8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
