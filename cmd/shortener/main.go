package main

import (
	"fmt"
	"github.com/kuznet1/urlshrt/internal/handler"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"net/http"
)

func main() {
	svc := service.NewService(&repository.MemoryRepo{})
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, h.Root)

	fmt.Println("Shortener service is starting at :8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
