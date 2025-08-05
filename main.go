package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nem0z/bitcoin-crawler/crawler"
	"github.com/nem0z/bitcoin-crawler/database"
	"github.com/nem0z/bitcoin-crawler/utils"
)

func main() {
	db, err := database.Init("./local.db")
	utils.Handle(err)

	addrs, err := utils.LoadAddr("./peers.json")
	utils.Handle(err)

	crawler, err := crawler.New(db, addrs...)
	utils.Handle(err)

	go crawler.LoadDB()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalCh
		fmt.Printf("Received signal: %v\n", sig)

		if err := crawler.Export("./export/nodes.json"); err != nil {
			log.Println("error exporting nodes.json", err)
		}

		crawler.SaveDB()

		os.Exit(0)
	}()

	select {}
}
