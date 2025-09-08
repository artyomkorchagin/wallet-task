package types_test

import (
	"encoding/json"
	"testing"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletUpdateRequest_UnmarshalJSON(t *testing.T) {
	jsonPayload := `{
        "valletId": "a1b2c3e4-5678-9012-3456-789012345678",
        "operationType": "DEPOSIT",
        "amount": 100
    }`

	var req types.WalletUpdateRequest
	err := json.Unmarshal([]byte(jsonPayload), &req)
	require.NoError(t, err)

	assert.Equal(t, "a1b2c3e4-5678-9012-3456-789012345678", req.WalletUUID)
	assert.Equal(t, "DEPOSIT", req.Operation)
	assert.Equal(t, 100, req.Amount)
}
func TestOperationTypeConstants(t *testing.T) {
	assert.Equal(t, "DEPOSIT", types.OperationTypeDeposit)
	assert.Equal(t, "WITHDRAW", types.OperationTypeWithdraw)
}

func TestNewWalletUpdateRequest(t *testing.T) {
	req := types.NewWalletUpdateRequest(
		"someuuid",
		"WITHDRAW",
		50,
		"ref-789",
	)

	assert.Equal(t, "someuuid", req.WalletUUID)
	assert.Equal(t, "WITHDRAW", req.Operation)
	assert.Equal(t, 50, req.Amount)
	assert.Equal(t, "ref-789", req.ReferenceID)
}
