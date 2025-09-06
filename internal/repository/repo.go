package repository

import (
	"context"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
)

type Repo interface {
	Put(ctx context.Context, url string) (model.URLID, error)
	Get(ctx context.Context, id model.URLID) (string, error)
	BatchPut(ctx context.Context, urls []string) ([]model.URLID, error)
	CreateUser(ctx context.Context) (int, error)
	UserUrls(ctx context.Context) (map[model.URLID]string, error)
	Ping(ctx context.Context) error
}

func NewRepo(cfg config.Config, logger *zap.Logger) (Repo, error) {
	if cfg.DatabaseDSN != "" {
		return NewDBRepo(cfg.DatabaseDSN, logger)
	}
	return NewMemoryRepo(cfg.FileStoragePath, logger)
}

type key int

const UserIDKey key = iota

func getUserID(ctx context.Context) (int, error) {
	val := ctx.Value(UserIDKey)
	id, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("unable to get id")
	}
	return id, nil
}
