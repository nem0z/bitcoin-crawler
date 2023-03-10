package handlers

import (
	"fmt"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/peer"
)

func Version(peer *peer.Peer, msg *message.Message) {
	fmt.Println("Deal with version", len(msg.Payload))
}

func Addr(peer *peer.Peer, msg *message.Message) {
	fmt.Println("Deal with version", len(msg.Payload))
}
