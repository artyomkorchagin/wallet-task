package router

import (
	"net/http"

	"github.com/artyomkorchagin/gin-boilerplate/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handlerFunc func(c *gin.Context) error

func (h *Handler) wrap(fn handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := fn(c)
		if err != nil {
			if httpErr, ok := err.(types.HTTPError); ok {
				h.logger.Error("error", zap.Error(httpErr))
				c.JSON(httpErr.Code, gin.H{"error": httpErr.Err.Error()})
			} else {
				h.logger.Error("error", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
	}
}
