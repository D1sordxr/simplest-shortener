package svc

import (
	"context"
	"math/rand"
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
	log pkg.Log
	storage
	dynamicRouter
}

func NewShortenerSvc(log pkg.Log, storage storage, dr dynamicRouter) *ShortenerSvc {
	return &ShortenerSvc{
		log:           log,
		storage:       storage,
		dynamicRouter: dr,
	}
}

func (s *ShortenerSvc) Create(ctx context.Context, url string) string {
	code := generateRandomCode()
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

const (
	newLength = 8
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomCode() string {
	newLen := rand.Intn(newLength) + 1

	code := make([]byte, newLen)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
