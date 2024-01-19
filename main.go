package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nem0z/bitcoin-crawler/crawler"
	"github.com/nem0z/bitcoin-crawler/peer"
	"github.com/nem0z/bitcoin-crawler/utils"
)

const peerIp = "190.64.134.52"
const peerPort = 8333

func main() {
	addr := &peer.Addr{Ip: peerIp, Port: peerPort}
	crawler, err := crawler.New(addr)
	utils.Handle(err)

	go crawler.Load("./export/nodes.json")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalCh
		fmt.Printf("Received signal: %v\n", sig)

		err = crawler.Export("./export/nodes.json")
		utils.Handle(err)

		os.Exit(0)
	}()

	select {}
}
