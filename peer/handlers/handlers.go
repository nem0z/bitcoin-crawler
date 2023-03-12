package handlers

import (
	"encoding/binary"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/peer"
)

func Version(ch chan *peer.Info) peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		version := int32(binary.LittleEndian.Uint32(msg.Payload[:4]))
		services := binary.LittleEndian.Uint64(msg.Payload[4:12])
		relay := true

		if version >= 70001 {
			relay = msg.Payload[len(msg.Payload)-1] != 0
		}

		p.Info = &peer.Info{Version: version, Services: services, Relay: relay}
		ch <- p.Info
	}
}

func Addr(ch chan []*peer.Addr) peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		ch <- ParseListAddr(msg.Payload)
	}
}
