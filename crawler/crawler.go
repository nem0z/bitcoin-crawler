package crawler

import (
	"fmt"
	"sync"

	chandlers "github.com/nem0z/bitcoin-crawler/crawler/handlers"
	"github.com/nem0z/bitcoin-crawler/peer"
	phandlers "github.com/nem0z/bitcoin-crawler/peer/handlers"
)

type Crawler struct {
	mu    sync.Mutex
	nodes map[string]*peer.Node
	out   chan *peer.Node
	addr  chan *peer.Addr
}

func New(addrs ...*peer.Addr) (*Crawler, error) {
	nodes := map[string]*peer.Node{}
	chOut := make(chan *peer.Node)
	chAddr := make(chan *peer.Addr)
	crawler := &Crawler{sync.Mutex{}, nodes, chOut, chAddr}

	for _, addr := range addrs {
		crawler.Add(addr)
	}

	go crawler.HandleResult()
	go crawler.HandleAddr(1000)

	go crawler.StartMonitoring()

	return crawler, nil
}

func (c *Crawler) add(addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nodes[addr] = nil
}

func (c *Crawler) Exist(addr string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.nodes[addr]
	return ok
}

func (c *Crawler) Add(addr *peer.Addr) {
	if c.Exist(addr.String()) {
		return
	}

	p, err := peer.New(addr.Ip, addr.Port, c.out)
	if err != nil {
		return
	}
	phandlers.DefaultRegister(p)
	p.Register("addr", chandlers.Addr(c.addr))

	err = p.Start()
	if err != nil {
		return
	}

	go Repeat(10, p.GetAddrNoError)

	go Delay(60, p.Close)
	c.add(addr.String())

}

func (c *Crawler) HandleResult() {
	for res := range c.out {
		addr := res.Addr.String()
		c.mu.Lock()
		if _, ok := c.nodes[addr]; ok {
			c.nodes[addr] = res
		}
		c.mu.Unlock()
	}
}

func (c *Crawler) HandleAddr(n int) {
	for i := 0; i < n; i++ {
		go func() {
			for addr := range c.addr {
				c.Add(addr)
			}
		}()
	}
}

func (c *Crawler) StartMonitoring() {
	go Repeat(10, c.Show)
}

func (c *Crawler) Show() {
	c.mu.Lock()
	defer c.mu.Unlock()

	cpt := len(c.nodes)
	cptProcessed := 0
	cptOk := 0

	for _, node := range c.nodes {
		if node != nil {
			cptProcessed++
			if node.Ping {
				cptOk++
			}
		}
	}

	fmt.Println()
	fmt.Println("Node discovered :", cpt)
	fmt.Println("Node processed :", cptProcessed)
	fmt.Println("Node ok :", cptOk)
	fmt.Println()
}
