package types_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestHTTPErrors_HaveCorrectCodes(t *testing.T) {
	err := errors.New("test")

	assert.Equal(t, http.StatusBadRequest, types.ErrBadRequest(err).Code)
	assert.Equal(t, http.StatusNotFound, types.ErrNotFound(err).Code)
	assert.Equal(t, http.StatusInternalServerError, types.ErrInternalServerError(err).Code)
	assert.Equal(t, http.StatusConflict, types.ErrConflict(err).Code)
}

func TestHTTPErrors_ErrorMethod(t *testing.T) {
	err := errors.New("something went wrong")
	httpErr := types.ErrBadRequest(err)

	assert.Equal(t, "something went wrong", httpErr.Error())
}

func TestErrorConstants(t *testing.T) {
	assert.Equal(t, "wallet not found", types.ErrWalletNotFound.Error())
	assert.Equal(t, "insufficient funds", types.ErrInsufficientFunds.Error())
	assert.Equal(t, "invalid operation type, must be DEPOSIT or WITHDRAW", types.ErrInvalidOperation.Error())
	assert.Equal(t, "concurrent update detected, retry required", types.ErrConcurrentUpdate.Error())
	assert.Equal(t, "operation with this reference_id already exists", types.ErrOperationExists.Error())
	assert.Equal(t, "database error", types.ErrDB.Error())
}
