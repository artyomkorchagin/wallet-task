package router

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_getBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid uuid - success", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/wallet/a1b2c3e4-5678-9012-3456-789012345678", nil)

		mockService.On("GetBalance", c, "a1b2c3e4-5678-9012-3456-789012345678").Return(500, nil)

		err := handler.getBalance(c)
		require.NoError(t, err)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"balance": 500}`, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/wallet/invalid-uuid", nil)

		mockService.On("GetBalance", c, "invalid-uuid").Return(500, nil)

		err := handler.getBalance(c)
		require.Error(t, err)
		httpErr, ok := err.(types.HTTPError)
		require.True(t, ok)
		assert.Equal(t, 400, httpErr.Code)
		assert.Contains(t, err.Error(), "walletUUID is not valid")
	})

	t.Run("service returns not found", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Params = gin.Params{{Key: "uuid", Value: "a1b2c3e4"}}

		mockService.On("GetBalance", mock.Anything, mock.Anything).Return(0, types.ErrWalletNotFound)

		err := handler.getBalance(c)
		require.Error(t, err)
		httpErr, ok := err.(types.HTTPError)
		require.True(t, ok)
		assert.Equal(t, 404, httpErr.Code)
		assert.Contains(t, err.Error(), "wallet not found")
	})
}

func TestHandler_updateBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid request - success", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		body := `{"walletId": "a1b2c3e4-5678-9012-3456-789012345678", "operationType": "DEPOSIT", "amount": 100}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/wallet/", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("UpdateBalance", mock.Anything, mock.MatchedBy(func(req *types.WalletUpdateRequest) bool {
			return req.WalletUUID == "a1b2c3e4-5678-9012-3456-789012345678" &&
				req.Operation == "DEPOSIT" &&
				req.Amount == 100 &&
				req.ReferenceID != ""
		})).Return(nil)

		err := handler.updateBalance(c)
		require.NoError(t, err)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"message": "balance updated"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/wallet/", strings.NewReader(`{invalid}`))
		c.Request.Header.Set("Content-Type", "application/json")

		err := handler.updateBalance(c)
		require.Error(t, err)
		httpErr, ok := err.(types.HTTPError)
		require.True(t, ok)
		assert.Equal(t, 400, httpErr.Code)
		assert.Contains(t, err.Error(), "invalid request body")
	})

	t.Run("missing required field", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		body := `{"walletId": "", "operationType": "DEPOSIT", "amount": 100}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/wallet/", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		err := handler.updateBalance(c)
		require.Error(t, err)
		httpErr, ok := err.(types.HTTPError)
		require.True(t, ok)
		assert.Equal(t, 400, httpErr.Code)
		assert.Contains(t, err.Error(), "walletUUID is empty")
	})

	t.Run("service returns insufficient funds", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewHandler(mockService, nil)

		body := `{"walletId": "a1b2c3e4.", "operationType": "WITHDRAW", "amount": 100}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/wallet/", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("UpdateBalance", mock.Anything, mock.Anything).Return(types.ErrInsufficientFunds)

		err := handler.updateBalance(c)
		require.Error(t, err)
		httpErr, ok := err.(types.HTTPError)
		require.True(t, ok)
		assert.Equal(t, 400, httpErr.Code)
		assert.Contains(t, err.Error(), "insufficient funds")
	})
}
