package repository

import (
	"encoding/json"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"net/http"
	"os"
	"sync"
)

type MemoryRepo struct {
	mutex sync.RWMutex
	store []string
	fname string
}

func NewMemoryRepo(fname string) (*MemoryRepo, error) {
	res := &MemoryRepo{fname: fname}
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
