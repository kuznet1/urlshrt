package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

var errNoDB = errors.New("database is not used")

type link struct {
	URL    string `json:"url"`
	UserID int    `json:"userId"`
}

type MemoryRepo struct {
	mutex      sync.RWMutex
	Store      []link `json:"store"`
	UsersCount int    `json:"usersCount"`
	fname      string
}

func NewMemoryRepo(fname string, logger *zap.Logger) (*MemoryRepo, error) {
	res := &MemoryRepo{fname: fname}

	if fname == "" {
		logger.Info("file storage path is empty, saving to file is disabled")
		return res, nil
	}

	_, err := os.Stat(fname)
	if err != nil {
		return res, nil
	}

	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open saved urls file %s: %w", fname, err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to read saved urls from file %s: %w", fname, err)
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
	userId, err := GetUserId(ctx)
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

	m.Store = append(m.Store, link{URL: url, UserID: userId})

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

	return m.Store[intID].URL, nil
}

func (m *MemoryRepo) BatchPut(ctx context.Context, urls []string) ([]model.URLID, error) {
	userId, err := GetUserId(ctx)
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

		m.Store = append(m.Store, link{URL: url, UserID: userId})
		res = append(res, model.URLID(len(m.Store)-1))
	}

	err1 := m.dump()
	if err1 != nil {
		return nil, err
	}

	return res, err
}

func (m *MemoryRepo) Ping(__ context.Context) error {
	return errNoDB
}

func (m *MemoryRepo) UserUrls(ctx context.Context) (map[model.URLID]string, error) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	res := make(map[model.URLID]string)
	for i, v := range m.Store {
		if v.UserID == userId {
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
