package peer

import (
	"io"
	"log"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/utils"
)

type Handler func(peer *Peer, message *message.Message)
type Handlers map[string]Handler

func (peer *Peer) Register(command string, handler Handler) {
	peer.handlers[command] = handler
}

func (peer *Peer) Handle() {
	for {
		msg, err := peer.Read()
		if err == io.EOF {
			continue
		}

		if err != nil {
			peer.Close()
			break
		}

		command := string(utils.RemoveTrailingZeros(msg.Command))
		handler, ok := peer.handlers[command]
		if ok {
			go handler(peer, msg)
			continue
		}

		log.Printf("Ignored message of type %s", string(msg.Command))
	}
}
