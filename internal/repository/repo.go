package repository

import (
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
)

type Repo interface {
	Put(url string) (model.URLID, error)
	Get(id model.URLID) (string, error)
	Ping() error
}

func NewRepo(cfg config.Config, logger *zap.Logger) (Repo, error) {
	if cfg.DatabaseDSN != "" {
		return NewDBRepo(cfg.DatabaseDSN, logger)
	}
	return NewMemoryRepo(cfg.FileStoragePath, logger)
}
