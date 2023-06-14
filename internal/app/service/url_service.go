package service

import (
	"context"
	"time"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
)

const timeout = 3 * time.Second

type UrlService struct {
	DB       domain.IUrlStorage
	Generate domain.IGenerateLinkService
}

func NewUrlService(db domain.IUrlStorage, service domain.IGenerateLinkService) *UrlService {
	return &UrlService{
		DB:       db,
		Generate: service,
	}
}

func (s *UrlService) CreateUrl(ctx context.Context, long string) (string, error) {

	if !CheckLink(long) {
		return "", domain.ErrorInvalidLink
	}

	urldata := &domain.URLData{}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return "", domain.ErrorGenerateTimeout
		default:
			short, now := s.Generate.GenerateShortLink()

			if _, err := s.DB.GetUrl(ctx, short); err == domain.ErrorLinkNotFound {
				urldata = domain.NewURLData(short, long, now)
			} else if err != nil {
				return "", err
			}
		}

		if urldata.URLShort != "" {
			break
		}
	}

	if err := s.DB.AddUrl(ctx, *urldata); err != nil {
		return "", err
	}

	return urldata.URLShort, nil
}

func (s UrlService) GetUrl(ctx context.Context, shortUrl string) (string, error) {

	if len(shortUrl) > 10 {
		return "", domain.ErrorInvalidShort
	}

	data, err := s.DB.GetUrl(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	return data.LongURL, nil
}
