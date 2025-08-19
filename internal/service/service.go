package service

import (
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/model"
	"github.com/kuznet1/urlshrt/internal/repository"
)

type Service struct {
	repo repository.Repo
	cfg  config.Config
}

func NewService(repo repository.Repo, cfg config.Config) Service {
	return Service{repo: repo, cfg: cfg}
}

func (svc Service) Shorten(url string) (string, error) {
	urlid, err := svc.repo.Put(url)
	if err != nil {
		return "", err
	}
	return svc.cfg.ShortenerPrefix + "/" + urlid.String(), nil
}

func (svc Service) Lengthen(id string) (string, error) {
	urlid, err := model.ParseURLID(id)
	if err != nil {
		return "", err
	}
	return svc.repo.Get(urlid)
}
