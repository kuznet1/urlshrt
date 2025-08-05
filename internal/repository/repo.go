package repository

import (
	"github.com/kuznet1/urlshrt/internal/model"
)

type Repo interface {
	Put(url string) (model.URLID, error)
	Get(id model.URLID) (string, error)
}
