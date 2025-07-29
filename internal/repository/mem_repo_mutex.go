package repository

import (
	"fmt"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"net/http"
	"sync"
)

type MemoryRepoMutex struct {
	mutex sync.RWMutex
	store []string
}

func (m *MemoryRepoMutex) Put(url string) (model.URLID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store = append(m.store, url)
	return model.URLID(len(m.store) - 1), nil
}

func (m *MemoryRepoMutex) Get(id model.URLID) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	intID := int(id.ID())
	if intID >= len(m.store) {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
	}

	return m.store[intID], nil
}
