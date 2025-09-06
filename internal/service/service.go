package service

import (
	"context"
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

func (svc Service) Shorten(ctx context.Context, url string) (string, error) {
	urlid, err := svc.repo.Put(ctx, url)
	return urlid.AsURL(svc.cfg.ShortenerPrefix), err
}

func (svc Service) BatchShorten(ctx context.Context, urls []string) ([]string, error) {
	if len(urls) == 0 {
		return []string{}, nil
	}

	urlids, err := svc.repo.BatchPut(ctx, urls)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, urlid := range urlids {
		res = append(res, urlid.AsURL(svc.cfg.ShortenerPrefix))
	}

	return res, nil
}

func (svc Service) UserUrls(ctx context.Context) ([]model.UrlsByUserResponseItem, error) {
	urlids, err := svc.repo.UserUrls(ctx)
	if err != nil {
		return nil, err
	}

	res := []model.UrlsByUserResponseItem{}
	for urlid, url := range urlids {
		res = append(res, model.UrlsByUserResponseItem{
			OriginalURL: url,
			ShortURL:    urlid.AsURL(svc.cfg.ShortenerPrefix),
		})
	}

	return res, nil
}

func (svc Service) Lengthen(ctx context.Context, id string) (string, error) {
	urlid, err := model.ParseURLID(id)
	if err != nil {
		return "", err
	}
	return svc.repo.Get(ctx, urlid)
}
