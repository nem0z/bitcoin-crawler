package handlers

import "github.com/nem0z/bitcoin-crawler/peer"

func DefaultRegister(p *peer.Peer) {
	p.Register("version", Version())
	p.Register("verack", Verack())
	p.Register("addr", Addr())
	p.Register("ping", Ping())
	p.Register("pong", Pong())
}
