package helpers

import (
	"net/http"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/gin-gonic/gin"
)

func HTTPError(c *gin.Context, err error) {
	switch err {
	case domain.ErrorGenerateTimeout:
		c.AbortWithError(http.StatusInternalServerError, err)
	case domain.ErrorLinkNotFound:
		c.AbortWithError(http.StatusNotFound, err)
	case domain.ErrorInvalidDecode:
		c.AbortWithError(http.StatusBadRequest, err)
	case domain.ErrorInvalidLink:
		c.AbortWithError(http.StatusBadRequest, err)
	default:
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
