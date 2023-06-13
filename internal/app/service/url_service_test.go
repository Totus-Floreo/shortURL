package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/Totus-Floreo/shortURL/internal/app/domain/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name        string
	Short       string
	Long        string
	AddedAt     int64
	GetUrlError error
	Error       error
}

var CreateTests = []TestCase{
	TestCase{
		Name:        "Success",
		Short:       "GoodLink12",
		Long:        "google.com",
		AddedAt:     1686557090,
		GetUrlError: domain.ErrorLinkNotFound,
		Error:       nil,
	},
	TestCase{
		Name:        "Timeout",
		Short:       "",
		Long:        "google.com",
		AddedAt:     0,
		GetUrlError: nil,
		Error:       domain.ErrorGenerateTimeout,
	},
	TestCase{
		Name:        "GetUrl Fatal",
		Short:       "",
		Long:        "google.com",
		AddedAt:     0,
		GetUrlError: ErrorDBShutdown,
		Error:       ErrorDBShutdown,
	},
	TestCase{
		Name:        "Add Url Error",
		Short:       "BadLuck123",
		Long:        "google.com",
		AddedAt:     1686557090,
		GetUrlError: domain.ErrorLinkNotFound,
		Error:       ErrorDBShutdown,
	},
}

var ErrorDBShutdown = errors.New("db shutdown")

func TestCreateUrl_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockIUrlStorage(ctrl)
	generator := mocks.NewMockIGenerateLinkService(ctrl)
	service := NewUrlService(db, generator)

	generator.EXPECT().GenerateShortLink().Return(CreateTests[0].Short, CreateTests[0].AddedAt)

	db.EXPECT().GetUrl(ctx, CreateTests[0].Short).Return(&domain.URLLong{}, CreateTests[0].GetUrlError)

	urldata := domain.NewURLData(CreateTests[0].Short, CreateTests[0].Long, CreateTests[0].AddedAt)

	db.EXPECT().AddUrl(ctx, *urldata).Return(nil)

	short, err := service.CreateUrl(ctx, CreateTests[0].Long)

	require.NoError(t, err)
	require.Equal(t, CreateTests[0].Short, short)
}

func TestCreateUrl_Timeout(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockIUrlStorage(ctrl)
	generator := mocks.NewMockIGenerateLinkService(ctrl)
	service := NewUrlService(db, generator)

	generator.EXPECT().GenerateShortLink().Return(CreateTests[1].Short, CreateTests[1].AddedAt).AnyTimes()

	db.EXPECT().GetUrl(ctx, CreateTests[1].Short).Return(&domain.URLLong{}, CreateTests[1].GetUrlError).AnyTimes()

	short, err := service.CreateUrl(ctx, CreateTests[1].Long)

	require.Error(t, err)
	require.Equal(t, CreateTests[1].Error, err)
	require.Equal(t, CreateTests[1].Short, short)
}

func TestCreateUrl_GetUrlError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockIUrlStorage(ctrl)
	generator := mocks.NewMockIGenerateLinkService(ctrl)
	service := NewUrlService(db, generator)

	generator.EXPECT().GenerateShortLink().Return(CreateTests[2].Short, CreateTests[2].AddedAt)

	db.EXPECT().GetUrl(ctx, CreateTests[2].Short).Return(&domain.URLLong{}, CreateTests[2].GetUrlError)

	short, err := service.CreateUrl(ctx, CreateTests[2].Long)

	require.Error(t, err)
	require.Equal(t, CreateTests[2].Error, err)
	require.Equal(t, CreateTests[2].Short, short)
}

func TestCreateUrl_AddUrlError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockIUrlStorage(ctrl)
	generator := mocks.NewMockIGenerateLinkService(ctrl)
	service := NewUrlService(db, generator)

	generator.EXPECT().GenerateShortLink().Return(CreateTests[3].Short, CreateTests[3].AddedAt)

	db.EXPECT().GetUrl(ctx, CreateTests[3].Short).Return(&domain.URLLong{}, CreateTests[3].GetUrlError)

	urldata := domain.NewURLData(CreateTests[3].Short, CreateTests[3].Long, CreateTests[3].AddedAt)

	db.EXPECT().AddUrl(ctx, *urldata).Return(ErrorDBShutdown)

	_, err := service.CreateUrl(ctx, CreateTests[3].Long)

	require.Error(t, err)
	require.Equal(t, CreateTests[3].Error, err)
}

func TestCreateUrl_InvalidLink(t *testing.T) {
	test := TestCase{
		Name:    "Invalid long link",
		Short:   "",
		Long:    "httpgooglecom",
		AddedAt: 0,
		Error:   domain.ErrorInvalidLink,
	}
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockIUrlStorage(ctrl)
	generator := NewGenerateLinkService()
	service := NewUrlService(db, generator)

	short, err := service.CreateUrl(ctx, test.Long)

	require.Equal(t, test.Error, err)
	require.Equal(t, test.Short, short)
}

func TestGetUrl(t *testing.T) {
	var tests = []TestCase{
		TestCase{
			Name:    "Success",
			Short:   "GoodLink12",
			Long:    "google.com",
			AddedAt: 1686557090,
			Error:   nil,
		},
		TestCase{
			Name:    "Fail",
			Short:   "BadLink123",
			Long:    "",
			AddedAt: 0,
			Error:   domain.ErrorLinkNotFound,
		},
	}
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockIUrlStorage(ctrl)
	generator := NewGenerateLinkService()
	service := NewUrlService(db, generator)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			urlLong := domain.URLLong{
				LongURL: test.Long,
				AddedAt: test.AddedAt,
			}

			db.EXPECT().GetUrl(ctx, test.Short).Return(&urlLong, test.Error)

			long, err := service.GetUrl(ctx, test.Short)

			require.Equal(t, test.Error, err)
			require.Equal(t, test.Long, long)
		})
	}
}
