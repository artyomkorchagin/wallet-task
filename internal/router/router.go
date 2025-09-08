package router

import (
	"net/http"

	walletservice "github.com/artyomkorchagin/wallet-task/internal/services/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	walletservice walletservice.ServiceInterface
	logger        *zap.Logger
}

func NewHandler(walletservice walletservice.ServiceInterface, logger *zap.Logger) *Handler {
	return &Handler{
		walletservice: walletservice,
		logger:        logger,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	main := router.Group("/")
	{
		main.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	apiv1 := router.Group("/api/v1/")
	{
		apiv1.GET("/wallet/:uuid", h.wrap(h.getBalance))
		apiv1.POST("/wallet", h.wrap(h.updateBalance))
	}
	h.logger.Info("Routes initialized")
	return router
}
