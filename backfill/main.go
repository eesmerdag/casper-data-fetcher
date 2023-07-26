package main

import (
	"cspr-fetcher/data"
	"cspr-fetcher/jobs"
	"database/sql"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	_ "github.com/lib/pq"
	"sync"
)

func main() {
	psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=casper-info sslmode=disable")
	db, err := sql.Open("postgres", psqlconn)
	if err = db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	rpcClient := sdk.NewRpcClient("http://54.180.220.20:7777/rpc")
	fetcher := jobs.NewBackfill(data.NewDBConnector(db), rpcClient)

	var wg sync.WaitGroup

	// numbers to be changed for future, currently about 2 million blocks
	for i := 1; i <= 1000; i++ {
		wg.Add(1)

		workerRange := 2000
		startLength := (i-1)*workerRange + 1
		endLength := i * workerRange

		go func() {
			defer wg.Done()
			fetcher.FetchDataWorker(uint64(startLength), uint64(endLength))
		}()
	}

	wg.Wait()
}
