package updater

import (
	"context"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg/db"
	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg/ethrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger    logrus.FieldLogger
	flags     Flags
	ethclient ethrpc.Client
	fromBlcok *big.Int
	toBlock   *big.Int
	dbClient  *db.DBClient
}

type ControllerParams struct {
	Logger    logrus.FieldLogger
	Flags     Flags
	Ethclient ethrpc.Client
	FromBlcok *big.Int
	ToBlock   *big.Int
	DbClient  *db.DBClient
}

func NewController(params ControllerParams) Controller {
	loggerWithFields := params.Logger.WithFields(logrus.Fields{
		"package": "updater",
		"struct":  "Controller",
	})
	controller := Controller{
		logger:    loggerWithFields,
		flags:     params.Flags,
		ethclient: params.Ethclient,
		fromBlcok: params.FromBlcok,
		toBlock:   params.ToBlock,
		dbClient:  params.DbClient,
	}
	return controller
}

func (c *Controller) Export(ctx context.Context) error {
	if err := c.runUntilSignal(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Controller) runUntilSignal(ctx context.Context) error {
	ctxWithCncl, cancel := context.WithCancel(ctx)
	defer cancel()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case s := <-sigc:
			c.logger.Info("got signal: ", s)
			cancel()
		case <-ctxWithCncl.Done():
		}

		signal.Stop(sigc)
		close(sigc)
	}()

	if err := c.run(ctxWithCncl); err != nil {
		c.logger.WithError(err).Error("done: ERROR")
		return err
	}

	return nil
}

func (c *Controller) run(ctx context.Context) error {
	transfers, err := c.ethclient.GetTransferWithBlock(ctx, c.fromBlcok, c.toBlock)
	if err != nil {
		return errors.Wrap(err, "failed to GetTransferWithBlock")
	}
	c.logger.Infof("extracted num of USDC transfer:%v", len(*transfers))

	res, err := c.dbClient.InsertUSDCTransfer(ctx, transfers)
	if err != nil {
		return errors.Wrap(err, "failed to GetTransferWithBlock")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get RowsAffected")
	}
	c.logger.Infof("inserted num of USDC transfer:%v", rowsAffected)

	return nil
}
