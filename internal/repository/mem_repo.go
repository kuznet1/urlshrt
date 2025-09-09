package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
	"net/http"
	"os"
	"slices"
	"sync"
)

var errNoDB = errors.New("database is not used")

type MemoryRepo struct {
	mutex sync.RWMutex
	store []string
	fname string
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

	err = json.NewDecoder(file).Decode(&res.store)
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

	return json.NewEncoder(file).Encode(m.store)
}

func (m *MemoryRepo) Put(url string) (model.URLID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	idx := slices.Index(m.store, url)
	if idx >= 0 {
		return model.URLID(idx), errs.NewDuplicatedURLError(url)
	}

	m.store = append(m.store, url)

	err := m.dump()
	if err != nil {
		return 0, err
	}

	return model.URLID(len(m.store) - 1), nil
}

func (m *MemoryRepo) Get(id model.URLID) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	intID := int(id.ID())
	if intID >= len(m.store) {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
	}

	return m.store[intID], nil
}

func (m *MemoryRepo) BatchPut(urls []string) ([]model.URLID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error
	var res []model.URLID
	for _, url := range urls {
		idx := slices.Index(m.store, url)
		if idx >= 0 {
			res = append(res, model.URLID(idx))
			err = errors.Join(err, errs.NewDuplicatedURLError(url))
			continue
		}
		m.store = append(m.store, url)
		res = append(res, model.URLID(len(m.store)-1))
	}

	err1 := m.dump()
	if err1 != nil {
		return nil, err
	}

	return res, err
}

func (m *MemoryRepo) Ping() error {
	return errNoDB
}
