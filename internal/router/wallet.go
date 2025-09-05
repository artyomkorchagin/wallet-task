package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getBalance(c *gin.Context) error {
	walletUUID := c.Param("uuid")
	fmt.Println(walletUUID)
	return nil
}

func (h *Handler) updateBalance(c *gin.Context) error {
	return nil
}
