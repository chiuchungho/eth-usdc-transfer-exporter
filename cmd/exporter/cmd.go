package updater

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg"
	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg/db"
	"github.com/chiuchungho/eth-usdc-transfer-exporter/pkg/ethrpc"
	"github.com/mitchellh/cli"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Command struct {
	flagSet  *flag.FlagSet
	defaults string
	flags    Flags
}

const (
	synopsis = "run exporter to pull data from ETH client and insert it to sqlite"
)

func NewCommand() (cli.Command, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags := Flags{
		SQLLitePath: f.String("sql-lite-path", "example.db", "local path to sql lite"),
		ETHRPCURL:   f.String("eth-rpc-url", "https://cloudflare-eth.com/", "ethereum rpc url"),
		FromBlcok:   f.String("from-block", "20369703", "range of first eth block"),
		ToBlock:     f.String("to-block", "20369703", "range of last eth block"),
		LogLevel:    f.String("log-level", "info", "logging level (debug, info, warn or error)"),
	}
	return &Command{
		flagSet:  f,
		defaults: pkg.GetFlagDefaults(f),
		flags:    flags,
	}, nil
}

func (c *Command) Run([]string) int {
	if err := c.flagSet.Parse(os.Args[2:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed to parse flags"))
		return 1
	}

	logger := logrus.New()
	logLevel, err := logrus.ParseLevel(*c.flags.LogLevel)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed to instantiate logger"))
		return 1
	}
	logger.SetLevel(logLevel)

	ctx := context.Background()
	ethClient, err := ethrpc.NewClient(ctx, ethrpc.RPCClientParams{
		RPCURL: *c.flags.ETHRPCURL,
		Logger: logger,
	})
	if err != nil {
		logger.Errorln(errors.Wrap(err, "failed to connect to eth rpc client"))
		return 1
	}

	fromBlcok := new(big.Int)
	fromBlcok, ok := fromBlcok.SetString(*c.flags.FromBlcok, 10)
	if !ok {
		logger.Errorln(errors.Wrapf(err, "failed to parse FromBlcok: %v", *c.flags.FromBlcok))
		return 1
	}
	toBlcok := new(big.Int)
	toBlcok, ok = toBlcok.SetString(*c.flags.ToBlock, 10)
	if !ok {
		logger.Errorln(errors.Wrapf(err, "failed to parse ToBlock: %v", *c.flags.ToBlock))
		return 1
	}
	logger.Infoln("block range:", fromBlcok, "-", toBlcok)

	logger.Infoln("sql lite path:", *c.flags.SQLLitePath)
	dbClient, err := db.NewDBClient(*c.flags.SQLLitePath, logger)
	if !ok {
		logger.Errorln(errors.Wrap(err, "failed to new db"))
		return 1
	}

	initTime := time.Now().UTC()
	logger.Infoln("initTime:", fmt.Sprint(initTime.UTC().Unix()))

	controller := NewController(ControllerParams{
		Ethclient: *ethClient,
		FromBlcok: fromBlcok,
		ToBlock:   toBlcok,
		DbClient:  dbClient,
		Logger:    logger,
		Flags:     c.flags,
	})

	if err := controller.Export(ctx); err != nil {
		logger.Errorln(errors.Wrap(err, "updater command failed"))
		return 1
	}
	logger.Infoln("done: OK - duration:", time.Since(initTime))
	return 0
}

func (c *Command) Help() string {
	return c.defaults
}

func (c *Command) Synopsis() string {
	return synopsis
}
