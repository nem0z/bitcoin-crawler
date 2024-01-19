package handlers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/message/payload"
	"github.com/nem0z/bitcoin-crawler/peer"
)

func Version() peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		version := int32(binary.LittleEndian.Uint32(msg.Payload[:4]))
		services := binary.LittleEndian.Uint64(msg.Payload[4:12])
		relay := true

		if version >= 70001 {
			relay = msg.Payload[len(msg.Payload)-1] != 0
		}

		p.Info = &peer.Info{Version: version, Services: services, Relay: relay}
	}
}

func Verack() peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		go p.ConsumeQueue()
	}
}

func Addr() peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		if msg.IsValid() {
			addrs := payload.ParseAddr(msg.Payload)
			for _, addr := range addrs {
				fmt.Println(addr)
			}
		}
	}
}

func Ping() peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		p.Pong(msg.Payload)
	}
}

func Pong() peer.Handler {
	return func(p *peer.Peer, msg *message.Message) {
		if bytes.Equal(msg.Payload, p.PingNonce) {
			p.PongAt = time.Now()
		}
	}
}
