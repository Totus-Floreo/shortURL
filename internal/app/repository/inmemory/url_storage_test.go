package inmemory

import (
	"context"
	"fmt"
	"testing"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	ID      int
	Short   string
	Long    string
	AddedAt int64
	Error   error
}

func TestAddUrl_Success(t *testing.T) {

	var test = TestCase{
		ID:      0,
		Short:   "S0mE__Lin4",
		Long:    "example.com",
		AddedAt: 1686557090,
		Error:   nil,
	}

	ctx := context.Background()

	urlStorage := NewUrlStorage()

	urldata := domain.NewURLData(test.Short, test.Long, test.AddedAt)
	err := urlStorage.AddUrl(ctx, *urldata)

	require.NoError(t, err)
}

func TestGetUrl(t *testing.T) {
	var tests = []TestCase{
		TestCase{
			ID:      1,
			Short:   "NormalLink",
			Long:    "example.com",
			AddedAt: 1686557090,
			Error:   nil,
		},
		TestCase{
			ID:      2,
			Short:   "BadLink123",
			Long:    "",
			AddedAt: 0,
			Error:   domain.ErrorLinkNotFound,
		},
	}

	ctx := context.Background()

	urlStorage := NewUrlStorage()
	urldata := domain.NewURLData(tests[0].Short, tests[0].Long, tests[0].AddedAt)
	_ = urlStorage.AddUrl(ctx, *urldata)

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test %v", test.ID), func(t *testing.T) {
			long, err := urlStorage.GetUrl(ctx, test.Short)

			require.Equal(t, test.Error, err)
			require.Equal(t, test.Long, long.LongURL)
			require.Equal(t, test.AddedAt, long.AddedAt)
		})
	}
}
