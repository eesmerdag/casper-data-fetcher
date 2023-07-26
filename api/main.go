package main

import (
	"cspr-fetcher/api/router"
	"cspr-fetcher/data"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=casper-info sslmode=disable")
	db, err := sql.Open("postgres", psqlconn)
	if err = db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	rt, err := router.NewRouter(data.NewDBConnector(db))
	if err != nil {
		log.Fatalf("error initializing router: %v", err)
	}
	fmt.Println("Service is ready...")

	if err != http.ListenAndServe(":1903", rt) {
		log.Fatalf("error initializing router: %v", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}
