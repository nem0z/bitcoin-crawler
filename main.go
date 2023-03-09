package main

import (
	"log"

	"github.com/nem0z/bitcoin-crawler/peer"
	"github.com/nem0z/bitcoin-crawler/utils"
)

func main() {
	//peerIp := "2a02:8108:8ac0:207b:d250:99ff:fe9e:792a"
	peerIp := "2a01:04f8:0201:2383:0000:0000:0000:0002"
	peerPort := 8333

	peer, err := peer.Create(peerIp, peerPort)
	utils.Handle(err)

	err = peer.Init()
	utils.Handle(err)

	err = peer.GetAddr()
	utils.Handle(err)

	err = peer.Mempool()
	utils.Handle(err)

	for i := 0; i < 15; i++ {
		msg, err := peer.Read()
		if err != nil {
			log.Fatal(err)
			break
		}
		msg.Display()
	}

}
