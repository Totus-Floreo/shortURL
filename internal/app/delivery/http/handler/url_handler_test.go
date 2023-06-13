package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/Totus-Floreo/shortURL/internal/app/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name         string
	Short        string
	Long         string
	AddedAt      int64
	ServiceError error
	StatusCode   int
}

var Tests = []TestCase{
	// Create
	TestCase{
		Name:         "Success",
		Short:        "GoodLink12",
		Long:         "google.com",
		ServiceError: nil,
		AddedAt:      1686557090,
		StatusCode:   http.StatusOK,
	},
	TestCase{
		Name:         "Bad Json",
		Short:        "",
		Long:         "",
		ServiceError: nil,
		AddedAt:      0,
		StatusCode:   http.StatusBadRequest,
	},
	TestCase{
		Name:         "Service Error",
		Short:        "",
		Long:         "google.com",
		ServiceError: domain.ErrorGenerateTimeout,
		AddedAt:      0,
		StatusCode:   http.StatusInternalServerError,
	},
	// Get
	TestCase{
		Name:         "Success",
		Short:        "GoodLink12",
		Long:         "google.com",
		ServiceError: nil,
		AddedAt:      1686557090,
		StatusCode:   http.StatusFound,
	},
	TestCase{
		Name:         "Service Error",
		Short:        "GoodLink12",
		Long:         "",
		ServiceError: ErrorDBShutdown,
		AddedAt:      0,
		StatusCode:   http.StatusInternalServerError,
	},
}

var ErrorDBShutdown = errors.New("db is shutdown")

// Create

func TestCreateUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockIUrlService(ctrl)
	handler := NewUrlHandler(service)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	service.EXPECT().CreateUrl(gomock.Any(), Tests[0].Long).Return(Tests[0].Short, Tests[0].ServiceError)

	router.POST("/", handler.CreateUrl)

	req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{"link": "%s"}`, Tests[0].Long)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	require.Equal(t, Tests[0].StatusCode, w.Code)
	require.Equal(t, Tests[0].Short, strings.Trim(w.Body.String(), `"/`))
}

func TestCreateUrl_BadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockIUrlService(ctrl)
	handler := NewUrlHandler(service)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.POST("/", handler.CreateUrl)

	req, _ := http.NewRequest("POST", "/", strings.NewReader(``))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	require.Equal(t, Tests[1].StatusCode, w.Code)
}

func TestCreateUrl_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockIUrlService(ctrl)
	handler := NewUrlHandler(service)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	service.EXPECT().CreateUrl(gomock.Any(), Tests[2].Long).Return(Tests[2].Short, Tests[2].ServiceError)

	router.POST("/", handler.CreateUrl)

	req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{"link": "%s"}`, Tests[2].Long)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	require.Equal(t, Tests[2].StatusCode, w.Code)
	require.Equal(t, Tests[2].Short, strings.Trim(w.Body.String(), `"/`))
}

// Get

func TestGetUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockIUrlService(ctrl)
	handler := NewUrlHandler(service)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	service.EXPECT().GetUrl(gomock.Any(), Tests[3].Short).Return(Tests[3].Long, Tests[3].ServiceError)

	router.GET("/:link", handler.GetUrl)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", Tests[3].Short), nil)

	router.ServeHTTP(w, req)

	require.Equal(t, Tests[3].StatusCode, w.Code)
	require.Equal(t, Tests[3].Long, strings.Trim(w.Body.String(), `"/`))
}

func TestGetUrl_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockIUrlService(ctrl)
	handler := NewUrlHandler(service)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	service.EXPECT().GetUrl(gomock.Any(), Tests[3].Short).Return(Tests[3].Long, Tests[3].ServiceError)

	router.GET("/:link", handler.GetUrl)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", Tests[3].Short), nil)

	router.ServeHTTP(w, req)

	require.Equal(t, Tests[3].StatusCode, w.Code)
}
