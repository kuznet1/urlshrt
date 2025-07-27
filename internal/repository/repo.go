package repository

import (
	"errors"
	"github.com/kuznet1/urlshrt/internal/model"
)

var ErrNotFound = errors.New("not found")

type Repo interface {
	Put(url string) (model.URLID, error)
	Get(id model.URLID) (string, error)
}
