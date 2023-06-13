package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/Totus-Floreo/shortURL/internal/app/domain/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestConnPool_BeginSuccess(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockpool := mocks.NewMockIPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockpool.EXPECT().Begin(ctx).Return(mockTx, nil)

	pool := Pool{mockpool}

	_, err := pool.Begin(ctx)

	require.NoError(t, err)
}

func TestConnPool_BeginError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockpool := mocks.NewMockIPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockpool.EXPECT().Begin(ctx).Return(mockTx, errors.New("pool error"))

	pool := Pool{mockpool}

	_, err := pool.Begin(ctx)

	require.Error(t, err)
}
