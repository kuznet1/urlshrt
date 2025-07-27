package repository

import (
	"fmt"
	"github.com/kuznet1/urlshrt/internal/model"
	"sync"
	"sync/atomic"
)

var _ Repo = (*MemoryRepo)(nil)

type MemoryRepo struct {
	counter atomic.Uint64
	store   sync.Map
}

func (m *MemoryRepo) Put(url string) (model.URLID, error) {
	id := m.counter.Add(1)
	m.store.Store(id, url)
	return model.URLID(id), nil
}

func (m *MemoryRepo) Get(id model.URLID) (string, error) {
	val, ok := m.store.Load(id.ID())
	if !ok {
		return "", ErrNotFound
	}

	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("url for shortening %s is not a string", id)
	}

	return s, nil
}
