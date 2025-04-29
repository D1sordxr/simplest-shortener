package svc

import (
	"context"
	"simplest-shortener/pkg"
)

type storage interface {
	Set(ctx context.Context, key, value string)
	Get(ctx context.Context, key string) (string, bool)
	Delete(ctx context.Context, key string)
}

type dynamicRouter interface {
	AddJob(mainURL, shortenedURL string)
}

type ShortenerSvc struct {
	storage
	dynamicRouter
	log pkg.Log
	gen pkg.Generator
}

func NewShortenerSvc(log pkg.Log, storage storage, dr dynamicRouter) *ShortenerSvc {
	return &ShortenerSvc{
		log:           log,
		storage:       storage,
		dynamicRouter: dr,
		gen:           new(pkg.SharedGenerator),
	}
}

const (
	newLengthEnv = 12
)

func (s *ShortenerSvc) Create(ctx context.Context, url string) string {
	if url == "" {
		s.log.Error("URL is empty")
		return ""
	}
	exist, ok := s.storage.Get(ctx, url)
	if ok {
		s.log.Info("URL already exists", "url", url, "code", exist)
		return exist
	}

	code := s.gen.GenerateRandomString(newLengthEnv)

	s.dynamicRouter.AddJob(url, code)

	s.storage.Set(ctx, code, url)

	s.log.Info("Created short URL", "code", code, "url", url)
	return code
}

func (s *ShortenerSvc) Get(ctx context.Context, code string) (string, bool) {
	url, ok := s.storage.Get(ctx, code)
	if !ok {
		s.log.Error("URL not found", "code", code)
		return "", false
	}
	s.log.Info("Retrieved URL", "code", code, "url", url)
	return url, true
}
