package main

import (
	"context"
	"errors"
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
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	log.Println("Build version: ", buildVersion)
	log.Println("Build date: ", buildDate)
	log.Println("Build commit: ", buildCommit)

	go func() {
		log.Println("pprof listening on :6060")
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

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

	svc := service.NewService(repo, cfg, logger)

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

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	connsClosed := make(chan any)
	go trapSignals(srv, connsClosed, logger)

	fmt.Println("Shortener service is starting at", cfg.ListenAddr)
	if cfg.EnableHTTPS {
		err = srv.ListenAndServeTLS(cfg.HTTPSCertFile, cfg.HTTPSCertKey)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-connsClosed
	fmt.Println("Shortener service is stopped")
}

func trapSignals(srv *http.Server, connsClosed chan any, logger *zap.Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-signals
	fmt.Println("Shortener service is stopping")
	err := srv.Shutdown(context.Background())
	if err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}
	close(connsClosed)
}
