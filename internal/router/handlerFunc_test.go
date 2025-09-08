package router

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func setupHandler(t *testing.T) (*Handler, *httptest.ResponseRecorder, *gin.Context) {
	logger := zaptest.NewLogger(t)
	mockService := new(MockWalletService)
	handler := NewHandler(mockService, logger)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return handler, w, c
}

func TestWrap_HTTPError(t *testing.T) {
	handler, w, c := setupHandler(t)

	httpErr := types.ErrNotFound(fmt.Errorf("not found"))
	ginHandler := handler.wrap(func(c *gin.Context) error {
		return httpErr
	})

	ginHandler(c)

	assert.Equal(t, 404, w.Code)
	assert.JSONEq(t, `{"error": "not found"}`, w.Body.String())
}

func TestWrap_GenericError(t *testing.T) {
	handler, w, c := setupHandler(t)

	genericErr := assert.AnError

	ginHandler := handler.wrap(func(c *gin.Context) error {
		return genericErr
	})

	ginHandler(c)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error": "assert.AnError general error for testing"}`, w.Body.String())
}

func TestWrap_NoError(t *testing.T) {
	handler, w, c := setupHandler(t)

	ginHandler := handler.wrap(func(c *gin.Context) error {
		return nil
	})

	ginHandler(c)

	assert.Equal(t, 200, w.Code)
	assert.Empty(t, w.Body.String())
}
