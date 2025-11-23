package repository

import (
	"context"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
)

// Repo abstracts storage for short URLs.
// Implementations must be safe for concurrent use where applicable and enforce per-user ownership.
// Methods: Put/Get single URL, BatchPut, BatchDelete, URLsByUser and user management helpers.
type Repo interface {
	Put(ctx context.Context, url string) (model.URLID, error)
	Get(ctx context.Context, id model.URLID) (string, error)
	BatchPut(ctx context.Context, urls []string) ([]model.URLID, error)
	CreateUser(ctx context.Context) (int, error)
	UserUrls(ctx context.Context) (map[model.URLID]string, error)
	BatchDelete(ctx context.Context, urlids []model.URLID) error
	Ping(ctx context.Context) error
}

// NewRepo performs a public package operation. Top-level handler/function.
func NewRepo(cfg config.Config, logger *zap.Logger) (Repo, error) {
	if cfg.DatabaseDSN != "" {
		return NewDBRepo(cfg, logger)
	}
	return NewMemoryRepo(cfg, logger)
}

type key int

// UserIDKey is the context key used to store the authenticated user's id.
const UserIDKey key = iota

// GetUserID extracts the authenticated user id from context.
// It returns an error if the id is missing or has an unexpected type.
func GetUserID(ctx context.Context) (int, error) {
	val := ctx.Value(UserIDKey)
	id, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("unable to get id")
	}
	return id, nil
}
