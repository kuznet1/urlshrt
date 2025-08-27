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
	return svc.cfg.ShortenerPrefix + "/" + urlid.String(), err
}

func (svc Service) BatchShorten(urls []string) ([]string, error) {
	if len(urls) == 0 {
		return []string{}, nil
	}

	urlids, err := svc.repo.BatchPut(urls)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, urlid := range urlids {
		res = append(res, svc.cfg.ShortenerPrefix+"/"+urlid.String())
	}

	return res, nil
}

func (svc Service) Lengthen(id string) (string, error) {
	urlid, err := model.ParseURLID(id)
	if err != nil {
		return "", err
	}
	return svc.repo.Get(urlid)
}
