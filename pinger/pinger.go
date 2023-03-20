package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/nem0z/bitcoin-crawler/peer"
	"github.com/nem0z/bitcoin-crawler/peer/handlers"
	"github.com/nem0z/bitcoin-crawler/utils"
)

type Peer struct {
	addr *peer.Addr
	ok   bool
}

func worker(addrs <-chan *peer.Addr, result chan *Peer) {
	for addr := range addrs {
		chStop := make(chan bool, 1)
		chPong := make(chan []byte, 1)
		ping, err := utils.CreateNonce(8)
		utils.Handle(err)

		go func() {
			ok := <-chStop
			result <- &Peer{addr, ok}
		}()

		go func() {
			pong := <-chPong
			chStop <- bytes.Equal(ping, pong)
		}()

		go func() {
			time.Sleep(time.Second * 5)
			chStop <- false
		}()

		p, err := peer.Create(addr.Ip, addr.Port, chStop)
		if err != nil {
			chStop <- false
			continue
		}

		p.Register("pong", handlers.Pong(chPong))
		err = p.Init()
		if err != nil {
			chStop <- false
			continue
		}

		err = p.Ping(ping)
		if err != nil {
			chStop <- true
			continue
		}
	}
}

func main() {
	peers, err := utils.ImportMap("./export/peers.json")
	utils.Handle(err)

	jobs := make(chan *peer.Addr)
	results := make(chan *Peer, len(peers))

	for i := 0; i < 100; i++ {
		go worker(jobs, results)
	}

	peersCpt := 0
	finised := false
	go func(finised *bool) {
		for ip, ok := range peers {
			if !ok {
				continue
			}
			peersCpt++
			jobs <- &peer.Addr{Ip: ip, Port: 8333}
		}
		*finised = true
	}(&finised)

	ok := 0
	processed := 0
	for i := 0; i < peersCpt || !finised; i++ {
		select {
		case result := <-results:
			processed++
			if result.ok {
				ok++
			}
			fmt.Printf("Valid / processed / total ( %v / %v / %v )\n", ok, processed, peersCpt)
		}
	}
}
