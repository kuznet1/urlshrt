package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/handler"
	"github.com/kuznet1/urlshrt/internal/middleware"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"github.com/kuznet1/urlshrt/internal/service/audit"
	"go.uber.org/zap"
	"log"
	"net/http"
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

	repo, err := repository.NewRepo(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewService(repo, cfg)

	if cfg.AuditFile != "" {
		listener, err := audit.NewFile(cfg.AuditFile)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		svc.Subscribe(listener)
	}

	if cfg.AuditURL != "" {
		svc.Subscribe(audit.NewURLAudit(cfg.AuditURL))
	}

	h := handler.NewHandler(svc, logger)
	requestLogger := middleware.NewRequestLogger(logger)
	auth := middleware.NewAuth(repo, cfg, logger)
	mux := chi.NewRouter()
	mux.Use(requestLogger.Logging, middleware.Compression, auth.Authentication)
	h.Register(mux)
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := repo.Ping(r.Context())
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
