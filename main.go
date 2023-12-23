package main

import (
	"time"

	"github.com/nem0z/bitcoin-crawler/peer"
	"github.com/nem0z/bitcoin-crawler/peer/handlers"
	"github.com/nem0z/bitcoin-crawler/utils"
)

const peerIp = "2604:a880:4:1d0::3d:2000"
const peerPort = 8333

func main() {
	p, err := peer.New(peerIp, peerPort)
	utils.Handle(err)

	handlers.DefaultRegister(p)

	err = p.Version()
	utils.Handle(err)

	err = p.Verack()
	utils.Handle(err)

	err = p.Ping()
	utils.Handle(err)

	err = p.GetAddr()
	utils.Handle(err)

	for {
		time.Sleep(3 * time.Second)
		p.Display()
	}
}
