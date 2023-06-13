package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/Totus-Floreo/shortURL/internal/app/domain/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	ID      int
	Short   string
	Long    string
	AddedAt int64
	Error   error
}

var (
	ErrBegin  error = errors.New("begin error")
	ErrScan   error = errors.New("scan error")
	ErrCommit error = errors.New("commit error")
	ErrExec   error = errors.New("exec error")
)

var Tests = []TestCase{
	//ADD
	TestCase{
		ID:      0,
		Short:   "S0mE__Lin4",
		Long:    "example.com",
		AddedAt: 1686557090,
		Error:   nil,
	},
	TestCase{
		ID:      1,
		Short:   "S0mE__Lin4",
		Long:    "example.com",
		AddedAt: 1686557090,
		Error:   ErrBegin,
	},
	TestCase{
		ID:      2,
		Short:   "S0mE__Lin4",
		Long:    "example.com",
		AddedAt: 1686557090,
		Error:   ErrExec,
	},
	TestCase{
		ID:      3,
		Short:   "S0mE__Lin4",
		Long:    "example.com",
		AddedAt: 1686557090,
		Error:   ErrCommit,
	},
	//GET
	TestCase{
		ID:      4,
		Short:   "S0mE__Lin4",
		Long:    "example.com",
		AddedAt: 1686557090,
		Error:   nil,
	},
	TestCase{
		ID:      5,
		Short:   "S0mE__Lin4",
		Long:    "",
		AddedAt: 0,
		Error:   ErrBegin,
	},
	TestCase{
		ID:      6,
		Short:   "сломанная-ссылка",
		Long:    "",
		AddedAt: 0,
		Error:   ErrScan,
	},
	TestCase{
		ID:      7,
		Short:   "S0mE__Lin4",
		Long:    "",
		AddedAt: 0,
		Error:   domain.ErrorLinkNotFound,
	},
	TestCase{
		ID:      8,
		Short:   "S0mE__Lin4",
		Long:    "",
		AddedAt: 0,
		Error:   ErrCommit,
	},
}

//AddUrl

func TestAddUrl_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	urldata := domain.NewURLData(Tests[0].Short, Tests[0].Long, Tests[0].AddedAt)
	err := urlStorage.AddUrl(ctx, *urldata)

	require.NoError(t, err)
}

func TestAddUrl_BeginError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	pool.EXPECT().Begin(gomock.Any()).Return(nil, ErrBegin)

	urldata := domain.NewURLData(Tests[1].Short, Tests[1].Long, Tests[1].AddedAt)
	err := urlStorage.AddUrl(ctx, *urldata)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[1].Error))
}

func TestAddUrl_ExecError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, ErrExec)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	urldata := domain.NewURLData(Tests[2].Short, Tests[2].Long, Tests[2].AddedAt)
	err := urlStorage.AddUrl(ctx, *urldata)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[2].Error))
}

func TestAddUrl_CommitError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	mockTx.EXPECT().Commit(gomock.Any()).Return(ErrCommit)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	urldata := domain.NewURLData(Tests[3].Short, Tests[3].Long, Tests[3].AddedAt)
	err := urlStorage.AddUrl(ctx, *urldata)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[3].Error))
}

// GetUrl

func TestGetUrl_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Tests[4].Short).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
		long := args[0].(*string)
		added := args[1].(*int64)
		*long = Tests[4].Long
		*added = Tests[4].AddedAt
		return nil
	})

	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	long, err := urlStorage.GetUrl(ctx, Tests[4].Short)

	require.NoError(t, err)
	require.Equal(t, Tests[4].Long, long.LongURL)
}

func TestGetUrl_BeginError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	pool.EXPECT().Begin(ctx).Return(nil, ErrBegin)

	_, err := urlStorage.GetUrl(ctx, Tests[5].Short)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[5].Error))
}

func TestGetUrl_ScanError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(ctx).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Tests[6].Short).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).Return(ErrScan)

	mockTx.EXPECT().Rollback(ctx).Return(nil)

	_, err := urlStorage.GetUrl(ctx, Tests[6].Short)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[6].Error))
}

func TestGetUrl_LinkNotFoundError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Tests[7].Short).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
		long := args[0].(*string)
		added := args[1].(*int64)
		*long = Tests[7].Long
		*added = Tests[7].AddedAt
		return pgx.ErrNoRows
	})

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	_, err := urlStorage.GetUrl(ctx, Tests[7].Short)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[7].Error))
}

func TestGetUrl_CommitError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	urlStorage := NewUrlStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Tests[8].Short).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
		long := args[0].(*string)
		added := args[1].(*int64)
		*long = Tests[8].Long
		*added = Tests[8].AddedAt
		return nil
	})

	mockTx.EXPECT().Commit(gomock.Any()).Return(ErrCommit)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	_, err := urlStorage.GetUrl(ctx, Tests[8].Short)

	require.Error(t, err)
	require.True(t, errors.Is(err, Tests[8].Error))
}
