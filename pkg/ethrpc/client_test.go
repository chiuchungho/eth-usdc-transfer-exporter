package ethrpc

import (
	"context"
	"math/big"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransferWithBlock(t *testing.T) {
	ctx := context.Background()
	cleint, err := NewClient(ctx, RPCClientParams{RPCURL: "https://cloudflare-eth.com/", Logger: logrus.New()})
	assert.NoErrorf(t, err, "expected to be no err to new eth rpc client")

	t.Run("Test_GetTransferWithBlock 20369703", func(t *testing.T) {
		transfers, err := cleint.GetTransferWithBlock(ctx, big.NewInt(20369703), big.NewInt(20369703))
		assert.NoErrorf(t, err, "expected to be no err to GetTransferWithBlock")
		assert.Equal(t, 35, len(*transfers))
	})
}
