package main

import (
	"cspr-fetcher/data"
	"cspr-fetcher/jobs"
	"database/sql"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
)

func main() {
	psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=casper-info sslmode=disable")
	db, err := sql.Open("postgres", psqlconn)
	if err = db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	fetcher := jobs.NewBlockInfoFetcher(data.NewDBConnector(db), sdk.NewRpcClient("http://54.180.220.20:7777/rpc"))

	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		fetcher.FetchBlockInfo()
	})

	c.Start()

	// dont stop fetcher
	select {}
}
