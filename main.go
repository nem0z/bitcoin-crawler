package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/nem0z/bitcoin-crawler/peer"
	"github.com/nem0z/bitcoin-crawler/peer/handlers"
	"github.com/nem0z/bitcoin-crawler/utils"
)

type PeerMap struct {
	arr map[string]bool
	mu  sync.Mutex
}

func (p *PeerMap) Contains(ip string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.arr[ip]
	return ok
}

func (p *PeerMap) Set(ip string, val bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.arr[ip] = val
}

func (p *PeerMap) Get(ip string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	val, ok := p.arr[ip]
	if !ok {
		return false, errors.New("Key not found in Peermap")
	}
	return val, nil
}

func (p *PeerMap) Export(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return utils.ExportMap(path, p.arr)
}

func main() {
	peerIp := "137.226.34.46"
	peerPort := 8333

	chVersion := make(chan *handlers.VersionOut)
	chAddr := make(chan []*peer.Addr)
	chCreate := make(chan *peer.Addr)
	stop := make(chan bool)

	p, err := peer.Create(peerIp, peerPort, nil)
	utils.Handle(err)

	p.Register("version", handlers.Version(chVersion))
	p.Register("addr", handlers.Addr(chAddr))

	err = p.Init()
	utils.Handle(err)

	err = p.GetAddr()
	utils.Handle(err)

	peers := &PeerMap{arr: make(map[string]bool)}

	for i := 0; i < 5000; i++ {
		go func(id int, addrs <-chan *peer.Addr) {
			for addr := range addrs {
				if peers.Contains(addr.Ip) {
					continue
				}

				chStop := make(chan bool)
				peers.Set(addr.Ip, false)

				p, err := peer.Create(addr.Ip, addr.Port, chStop)
				if err != nil {
					continue
				}

				p.Register("version", handlers.Version(chVersion))
				p.Register("addr", handlers.Addr(chAddr))
				err = p.Init()
				if err != nil {
					continue
				}

				err = p.GetAddr()
				if err != nil {
					continue
				}
			}
		}(i, chCreate)
	}

	cpt := 0

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("Found %v valid peers over %v performed\n", cpt, len(peers.arr))
			fmt.Println(runtime.NumGoroutine())
			if cpt >= 1000 {
				stop <- true
			}
		}
	}()

Loop:
	for {
		select {
		case version := <-chVersion:
			peers.Set(version.Addr.Ip, true)
			cpt++

		case addrs := <-chAddr:
			for _, addr := range addrs {
				chCreate <- addr
			}
		case <-stop:
			break Loop
		}
	}

	err = peers.Export("./export/peers.json")
	utils.Handle(err)
}
