package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type IPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}
