package postgresql

import (
	"context"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/jackc/pgx/v5"
)

type UrlStorage struct {
	Pool domain.IPool
}

func NewUrlStorage(pool domain.IPool) *UrlStorage {
	return &UrlStorage{
		Pool: pool,
	}
}

func (s *UrlStorage) AddUrl(ctx context.Context, urlData domain.URLData) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO links(short, long, added) VALUES ($1, $2, $3)", urlData.URLShort, urlData.LongURL, urlData.AddedAt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *UrlStorage) GetUrl(ctx context.Context, shortUrl string) (*domain.URLLong, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return &domain.URLLong{}, err
	}
	defer tx.Rollback(ctx)

	urllong := &domain.URLLong{}
	if err := tx.QueryRow(ctx, "SELECT long, added FROM links WHERE short = $1", shortUrl).Scan(&urllong.LongURL, &urllong.AddedAt); err != nil {
		if err == pgx.ErrNoRows {
			return &domain.URLLong{}, domain.ErrorLinkNotFound
		} else {
			return &domain.URLLong{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return &domain.URLLong{}, err
	}

	return urllong, nil
}
