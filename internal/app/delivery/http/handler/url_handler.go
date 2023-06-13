package http

import (
	"net/http"

	"github.com/Totus-Floreo/shortURL/internal/app/delivery/http/helpers"
	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	Service domain.IUrlService
}

func NewUrlHandler(service domain.IUrlService) *UrlHandler {
	return &UrlHandler{
		Service: service,
	}
}

func (h *UrlHandler) CreateUrl(c *gin.Context) {
	var long domain.URLLong
	if err := c.BindJSON(&long); err != nil {
		helpers.HTTPError(c, domain.ErrorInvalidDecode)
		return
	}

	short, err := h.Service.CreateUrl(c, long.LongURL)
	if err != nil {
		helpers.HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, short)
}

func (h *UrlHandler) GetUrl(c *gin.Context) {
	short := c.Param("link")

	long, err := h.Service.GetUrl(c, short)
	if err != nil {
		helpers.HTTPError(c, err)
		return
	}

	c.JSON(http.StatusFound, long)
}
