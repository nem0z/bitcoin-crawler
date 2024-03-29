package handlers

import (
	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/message/payload"
	"github.com/nem0z/bitcoin-crawler/peer"
)

func Addr(ch chan *peer.Addr) peer.Handler {
	cpt := 0
	return func(p *peer.Peer, msg *message.Message) {
		if msg.IsValid() {
			if cpt++; cpt >= 5 {
				p.Close()
			}
			addrs := payload.ParseAddr(msg.Payload)
			for _, addr := range addrs {
				ch <- &peer.Addr{Ip: addr.Ip, Port: int(addr.Port)}
			}
		}
	}
}
