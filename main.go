package main

import (
	"github.com/nem0z/bitcoin-crawler/crawler"
	"github.com/nem0z/bitcoin-crawler/utils"
)

const peerIp = "2604:a880:4:1d0::3d:2000"
const peerPort = 8333

func main() {
	_, err := crawler.New(peerIp, peerPort)
	utils.Handle(err)

	select {}
}
