package service

import (
	"context"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/model"
	"github.com/kuznet1/urlshrt/internal/repository"
	"go.uber.org/zap"
)

// Service contains the application business logic atop the storage layer.
// It coordinates repository operations and publishes audit events.
type Service struct {
	repo   repository.Repo
	cfg    config.Config
	logger *zap.Logger
	subs   []AuditSubscriber
}

// AuditSubscriber is notified about URL creation events.
// Implementations may forward events to files, HTTP endpoints, or external systems.
type AuditSubscriber interface {
	OnAuditEvt(ctx context.Context, userID int, action model.AuditAction, url string) error
}

// NewService constructs a Service with the given repository and configuration.
func NewService(repo repository.Repo, cfg config.Config, logger *zap.Logger) Service {
	return Service{repo: repo, cfg: cfg, logger: logger}
}

// Shorten validates and stores a single URL and returns its short identifier.
// If the URL already exists for the user, a DuplicatedURLError is returned.
func (svc *Service) Shorten(ctx context.Context, url string) (string, error) {
	urlid, err := svc.repo.Put(ctx, url)
	svc.fire(ctx, model.ActionShorten, url)
	return urlid.AsURL(svc.cfg.ShortenerPrefix), err
}

// BatchShorten stores multiple URLs at once and returns their identifiers in the same order.
// Errors for individual items are combined; duplicates are reported as DuplicatedURLError.
func (svc *Service) BatchShorten(ctx context.Context, urls []string) ([]string, error) {
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

// BatchDelete removes the given short ids that belong to the current user.
// The actual deletion strategy (immediate vs. batched) depends on the repository implementation.
func (svc *Service) BatchDelete(ctx context.Context, ids []string) error {
	var urlids []model.URLID
	for _, id := range ids {
		urlid, err := model.ParseURLID(id)
		if err != nil {
			return err
		}
		urlids = append(urlids, urlid)
	}

	return svc.repo.BatchDelete(ctx, urlids)
}

// UserUrls returns a list of the user's URLs along with their absolute short forms.
func (svc *Service) UserUrls(ctx context.Context) ([]model.UrlsByUserResponseItem, error) {
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

// Lengthen resolves a short identifier back to the original URL.
func (svc *Service) Lengthen(ctx context.Context, id string) (string, error) {
	urlid, err := model.ParseURLID(id)
	if err != nil {
		return "", err
	}

	url, err := svc.repo.Get(ctx, urlid)
	svc.fire(ctx, model.ActionFollow, url)
	return url, err
}

// Subscribe registers an AuditSubscriber that will be notified about URL creation events.
func (svc *Service) Subscribe(sub AuditSubscriber) {
	svc.subs = append(svc.subs, sub)
}

func (svc *Service) fire(ctx context.Context, action model.AuditAction, url string) {
	userID, err := repository.GetUserID(ctx)
	if err != nil {
		svc.logger.Error("audit event handling error", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(ctx, svc.cfg.AuditURLTimeout)
	defer cancel()
	for _, sub := range svc.subs {
		err = sub.OnAuditEvt(ctx, userID, action, url)
		if err != nil {
			svc.logger.Error("audit event handling error", zap.Error(err))
		}
	}
}
