package postgresql

import (
	"context"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/jackc/pgx/v5"
)

type Pool struct {
	connPool domain.IPool
}

func (pool *Pool) Begin(ctx context.Context) (pgx.Tx, error) {
	Tx, err := pool.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return Tx, nil
}
