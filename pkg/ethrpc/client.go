package ethrpc

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg/db"
	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg/erc20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	usdcAddress     = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	logTreansferSig = "Transfer(address,address,uint256)"
)

type Client struct {
	client ethclient.Client
	logger logrus.FieldLogger
}

type RPCClientParams struct {
	RPCURL string
	Logger logrus.FieldLogger
}

func NewClient(ctx context.Context, params RPCClientParams) (*Client, error) {
	loggerWithFields := params.Logger.WithFields(logrus.Fields{
		"package": "ethrpc",
		"struct":  "Client",
	})

	client, err := ethclient.DialContext(ctx, params.RPCURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial rpc url=%s", params.RPCURL)
	}
	return &Client{
		client: *client,
		logger: loggerWithFields}, nil
}

type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func (c *Client) GetTransferWithBlock(ctx context.Context, fromBlcok, toBlcok *big.Int) (*[]db.Transfer, error) {
	c.logger.Info("GetTransferWithBlock for USDC")
	contractAddress := common.HexToAddress(usdcAddress)
	query := ethereum.FilterQuery{
		FromBlock: fromBlcok,
		ToBlock:   toBlcok,
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := c.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to to FilterLogs")
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(erc20.Erc20ABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte(logTreansferSig)
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	out := make([]db.Transfer, 0)
	for _, vLog := range logs {
		if vLog.Topics[0].Hex() == logTransferSigHash.Hex() {
			var transferEvent LogTransfer

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				return nil, errors.Wrap(err, "failed to UnpackIntoInterface")
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			out = append(out, db.Transfer{
				From:  transferEvent.From.Hex(),
				To:    transferEvent.To.Hex(),
				Value: transferEvent.Value.String(),
			})
		}
	}

	c.logger.Infof("Extracted nums of transfer: %v", len(out))
	return &out, nil
}
