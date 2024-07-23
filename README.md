# eth-usdc-transfer-exporter

## Task
Write a go application that uses the eth json rpc api to extract all USDC token transfers included in a given ethereum mainnet block and saves them into a sqlite database (each row should include at least sender / recipient and the value that was transferred)

## How to run it
You run it with script or docker.
- by script: 

1. `cd cmd`
1. `chmod u+x run.sh`
1. `./run.sh`


- docker

`docker compose up && docker compose rm -f`

## Test
I have included 2 tests here only because of the time pressure. If it is for proper appication, I would also like to include integration test for the whole app.
1. pkg/ethrpc/client_test.go
1. pkg/db/insert_test.go

## System Design
It is an application that can be run with cli command. It also support the following arguments
```
1. ETHRPCURL (eth rpc client)
2. FROM_BLOCK (first block to select)
3. TO_BLOCK (last block to select)
4. LOG_LEVEL
5. SQL_LITE_PATH (path to local sql lite)
```

The availibiliy to input different args provides user to select different block range of USDC transfer.

## further improvement
1. This task is only designed to select 1 block of USDC transfer. It only has limited amount of USDC transfor for 1 block. If we want to select all the past records of USDC transfer, this application can be improved with adding go goroutine to `cmd/exporter/controller.go` -> `func (c *Controller) run(ctx context.Context) error`. By doing the `GetTransferWithBlock` and `InsertUSDCTransfer` at the same time, it can process the transfer records in batches. 
1. If the application needs to process all the ERC20 Transfer, I would build it as 2 app. 1 app as `extractor` which connect to eth rpc client and collect all the data. And then `extractor` pass the data with a message queue to another app as `exporter`. `exporter` will insert the data into big query or big table.
1. sqlite is not as fast as postgreSQL, I would use postgreSQL to replace the current setting
1. I assume that the data needs to be served to our FE or client, we can add an api server to serve the data.