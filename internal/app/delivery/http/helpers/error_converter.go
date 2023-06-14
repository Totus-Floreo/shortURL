package helpers

import (
	"net/http"

	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/gin-gonic/gin"
)

func HTTPError(c *gin.Context, err error) {
	ginError := gin.Error{
		Err: err,
	}
	switch err {
	case domain.ErrorGenerateTimeout:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ginError.JSON())
	case domain.ErrorLinkNotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, ginError.JSON())
	case domain.ErrorInvalidDecode:
		c.AbortWithStatusJSON(http.StatusBadRequest, ginError.JSON())
	case domain.ErrorInvalidShort:
		c.AbortWithStatusJSON(http.StatusBadRequest, ginError.JSON())
	case domain.ErrorInvalidLink:
		c.AbortWithStatusJSON(http.StatusBadRequest, ginError.JSON())
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ginError.JSON())
	}
}
