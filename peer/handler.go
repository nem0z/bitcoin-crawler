package peer

import (
	"io"

	"github.com/nem0z/bitcoin-crawler/message"
)

type Handler func(peer *Peer, message *message.Message)
type Handlers map[string]Handler

func (peer *Peer) Register(command string, handler Handler) {
	peer.handlers[command] = handler
}

func (peer *Peer) Handle() {
	for {
		if peer.conn == nil {
			return
		}

		msg, err := peer.Read()
		if err == io.EOF {
			continue
		}

		if err != nil {
			peer.Close()
			return
		}

		command := message.ResolveCommandName(msg.Command)
		handler, ok := peer.handlers[command]
		if ok && msg.IsValid() {
			// log.Println("Handle message :", message.ResolveCommandName(msg.Command))
			go handler(peer, msg)
			continue
		}
	}
}
