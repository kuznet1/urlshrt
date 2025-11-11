package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

var errNoDB = errors.New("database is not used")

type link struct {
	URL       string `json:"url"`
	UserID    int    `json:"userID"`
	IsDeleted bool   `json:"isDeleted"`
}

type MemoryRepo struct {
	batchRemover
	mutex      sync.RWMutex
	Store      []*link `json:"store"`
	UsersCount int     `json:"usersCount"`
	fname      string
	logger     *zap.Logger
}

func NewMemoryRepo(cfg config.Config, logger *zap.Logger) (*MemoryRepo, error) {
	res := &MemoryRepo{batchRemover: newBatchRemover(cfg), fname: cfg.FileStoragePath, logger: logger}
	go res.deletionWorker(res.deleteImpl)

	if cfg.FileStoragePath == "" {
		logger.Info("file storage path is empty, saving to file is disabled")
		return res, nil
	}

	_, err := os.Stat(cfg.FileStoragePath)
	if err != nil {
		return res, nil
	}

	file, err := os.Open(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open saved urls file %s: %w", cfg.FileStoragePath, err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to read saved urls from file %s: %w", cfg.FileStoragePath, err)
	}

	return res, nil
}

func (m *MemoryRepo) dump() error {
	if m.fname == "" {
		return nil
	}
	file, err := os.Create(m.fname)
	if err != nil {
		return fmt.Errorf("failed to save urls to file %s: %w", m.fname, err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(m)
}

func (m *MemoryRepo) Put(ctx context.Context, url string) (model.URLID, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return 0, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, v := range m.Store {
		if v.URL == url {
			return model.URLID(i), errs.NewDuplicatedURLError(url)
		}
	}

	m.Store = append(m.Store, &link{URL: url, UserID: userID})

	err = m.dump()
	if err != nil {
		return 0, err
	}

	return model.URLID(len(m.Store) - 1), nil
}

func (m *MemoryRepo) Get(_ context.Context, id model.URLID) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	intID := int(id.ID())
	if intID >= len(m.Store) {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
	}

	res := m.Store[intID]
	if res.IsDeleted {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q is deleted", id), http.StatusGone)
	}

	return res.URL, nil
}

func (m *MemoryRepo) BatchPut(ctx context.Context, urls []string) ([]model.URLID, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var res []model.URLID
	for _, url := range urls {
		for i, v := range m.Store {
			if v.URL == url {
				res = append(res, model.URLID(i))
				err = errors.Join(err, errs.NewDuplicatedURLError(url))
				continue
			}
		}

		m.Store = append(m.Store, &link{URL: url, UserID: userID})
		res = append(res, model.URLID(len(m.Store)-1))
	}

	err1 := m.dump()
	if err1 != nil {
		return nil, err
	}

	return res, err
}

func (m *MemoryRepo) deleteImpl(reqs []deleteLinkReq) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, req := range reqs {
		if int(req.urlid) > len(m.Store) {
			m.logger.Error("no such link")
			continue
		}

		l := m.Store[req.urlid]
		if l.UserID != req.userID {
			m.logger.Error("access denied")
			continue
		}
		l.IsDeleted = true
	}
}

func (m *MemoryRepo) Ping(__ context.Context) error {
	return errNoDB
}

func (m *MemoryRepo) UserUrls(ctx context.Context) (map[model.URLID]string, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	res := make(map[model.URLID]string)
	for i, v := range m.Store {
		if v.UserID == userID {
			res[model.URLID(i)] = v.URL
		}
	}

	return res, nil
}

func (m *MemoryRepo) CreateUser(_ context.Context) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	res := m.UsersCount
	m.UsersCount++
	return res, nil
}
