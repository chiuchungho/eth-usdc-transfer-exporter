package main

import (
	"log"
	"os"

	"github.com/chiuchungho/eth-usdc-transfer-exporter/cmd/exporter"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("eth-usdc-transfer", "")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"exporter": updater.NewCommand,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
