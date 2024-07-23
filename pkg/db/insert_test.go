package db

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_InsertUSDCTransfer(t *testing.T) {
	ctx := context.Background()
	cleint, err := NewDBClient("example.db", logrus.New())
	assert.NoErrorf(t, err, "expected to be no err to new eth rpc client")

	t.Run("Test_InsertUSDCTransfer", func(t *testing.T) {
		transfer := []Transfer{
			{
				From:  "1",
				To:    "2",
				Value: "3",
			},
		}
		res, err := cleint.InsertUSDCTransfer(ctx, &transfer)
		spew.Dump(res)
		assert.NoErrorf(t, err, "expected to be no err to InsertUSDCTransfer")
		rowsAffected, err := res.RowsAffected()
		assert.NoError(t, err, "expected to be no err for res.RowsAffected()")
		assert.Equal(t, int64(1), rowsAffected)
	})
}
