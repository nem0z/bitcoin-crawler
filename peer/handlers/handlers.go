package handlers

import (
	"encoding/binary"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/peer"
)

type VersionOut struct {
	Addr *peer.Addr
	Info *peer.Info
}

func Version(ch chan *VersionOut) peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		version := int32(binary.LittleEndian.Uint32(msg.Payload[:4]))
		services := binary.LittleEndian.Uint64(msg.Payload[4:12])
		relay := true

		if version >= 70001 {
			relay = msg.Payload[len(msg.Payload)-1] != 0
		}

		p.Info = &peer.Info{Version: version, Services: services, Relay: relay}

		versionOut := &VersionOut{p.SelfAddr(), p.Info}
		ch <- versionOut
	}
}

func Addr(ch chan []*peer.Addr) peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		ch <- ParseListAddr(msg.Payload)
		p.Close()
	}
}

func Pong(ch chan []byte) peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		ch <- msg.Payload
	}
}
