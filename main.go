package main

import (
	"fmt"

	"github.com/nem0z/bitcoin-crawler/peer"
	"github.com/nem0z/bitcoin-crawler/peer/handlers"
	"github.com/nem0z/bitcoin-crawler/utils"
)

func main() {
	peerIp := "2a00:8a60:e012:a00::21"
	peerPort := 8333

	chVersion := make(chan *peer.Info)
	chAddr := make(chan []*peer.Addr)

	p, err := peer.Create(peerIp, peerPort)
	utils.Handle(err)

	p.Register("version", handlers.Version(chVersion))
	p.Register("addr", handlers.Addr(chAddr))

	err = p.Init()
	utils.Handle(err)

	err = p.GetAddr()
	utils.Handle(err)

	for {
		select {

		case info := <-chVersion:
			fmt.Println("Version received :", info)

		case addrs := <-chAddr:
			fmt.Println("Addr received :", len(addrs))
		}
	}
}
