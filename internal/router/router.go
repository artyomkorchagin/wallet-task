package router

import (
	"net/http"

	//_ "github.com/artyomkorchagin/wallet-task/docs"

	walletservice "github.com/artyomkorchagin/wallet-task/internal/services/wallet"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	walletservice *walletservice.Service
	logger        *zap.Logger
}

func NewHandler(walletservice *walletservice.Service, logger *zap.Logger) *Handler {
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

		main.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiv1 := router.Group("/api/v1/")
	{
		apiv1.POST("/wallet", h.wrap(h.getBalance))
		apiv1.GET("/wallet/:uuid", h.wrap(h.updateBalance))
	}
	h.logger.Info("Routes initialized")
	return router
}
