package handler

import (
	"errors"
	"fmt"
	"net/http"

	"api-gateway/internal/handlergen"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	gateway service.GatewayService
}

func NewHandler(gateway service.GatewayService) *Handler {
	return &Handler{
		gateway: gateway,
	}
}

func (h *Handler) PostShorten(c *gin.Context) {
	var body handlergen.PostShortenJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, handlergen.ErrorResponse{Error: &errMsg})
	}

	shortUrl, err := h.gateway.ShortenURL(c, *body.LongUrl)
	if err != nil {
		handleError(c, err)
		return
	}

	var response handlergen.ShortenResponse
	response.ShortUrl = &shortUrl
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetShortUrl(c *gin.Context, shortUrl string) {
	originalURL, err := h.gateway.GetLongURL(c, shortUrl)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *Handler) PostStats(c *gin.Context) {
	// TODO
	panic("impelement me")
}

func handleError(c *gin.Context, err error) {
	var (
		response handlergen.ErrorResponse
		errMsg   string
	)
	switch {
	case errors.Is(err, service.ErrInvalidURL):
		errMsg = fmt.Sprintf("invalid URL provided: %s", err.Error())
		response.Error = &errMsg
		c.JSON(http.StatusBadRequest, response)
	case errors.Is(err, service.ErrURLNotFound):
		errMsg = fmt.Sprintf("URL not found: %s", err.Error())
		response.Error = &errMsg
		c.JSON(http.StatusNotFound, response)
	default:
		errMsg = fmt.Sprintf("internal error: %s", err.Error())
		response.Error = &errMsg
		c.JSON(http.StatusInternalServerError, response)
	}
}
