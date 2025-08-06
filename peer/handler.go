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
		select {
		case <-peer.ctx.Done():
			return

		default:
			if peer.conn == nil {
				peer.Close()
				return
			}

			msg, err := peer.Read()
			if err == io.EOF {
				continue
			}

			if err != nil || !msg.IsValid() {
				peer.Close()
				return
			}

			command := message.ResolveCommandName(msg.Command)
			if handler, ok := peer.handlers[command]; ok {
				// log.Println("Handle message :", message.ResolveCommandName(msg.Command))
				go handler(peer, msg)
			}
		}
	}
}
