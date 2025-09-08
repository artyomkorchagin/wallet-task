package router

import (
	"fmt"
	"net/http"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) getBalance(c *gin.Context) error {
	walletUUID := c.Param("uuid")

	balance, err := h.walletservice.GetBalance(c, walletUUID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
	return nil
}

func (h *Handler) updateBalance(c *gin.Context) error {
	var wur *types.WalletUpdateRequest

	if err := c.ShouldBindJSON(&wur); err != nil {
		return types.ErrBadRequest(fmt.Errorf("invalid request body: %w", err))
	}

	// it's supposed to be generated on the client side,
	// but for this task there isnt one, so i'll simulate it
	wur.ReferenceID = uuid.New().String()

	if err := h.walletservice.UpdateBalance(c, wur); err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "balance updated"})
	return nil
}
