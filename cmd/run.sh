#!/bin/bash
set -eo pipefail

go build

export ETHRPCURL=https://cloudflare-eth.com/
export FROM_BLOCK="20369703"
export TO_BLOCK="20369703"

go run main.go exporter