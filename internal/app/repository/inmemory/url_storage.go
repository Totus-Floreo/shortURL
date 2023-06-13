package inmemory

import (
	"context"
	"sync"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
)

type UrlStorage struct {
	Mux     *sync.RWMutex
	Storage map[string]domain.URLLong
}

func NewUrlStorage() *UrlStorage {
	return &UrlStorage{
		Mux:     new(sync.RWMutex),
		Storage: map[string]domain.URLLong{},
	}
}

func (s *UrlStorage) AddUrl(ctx context.Context, urlData domain.URLData) error {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	s.Storage[urlData.URLShort] = urlData.URLLong
	return nil
}

func (s *UrlStorage) GetUrl(ctx context.Context, shortUrl string) (*domain.URLLong, error) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	longUrl, ok := s.Storage[shortUrl]
	if !ok {
		return &domain.URLLong{}, domain.ErrorLinkNotFound
	}

	return &longUrl, nil
}
