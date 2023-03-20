package peer

import (
	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/message/payload"
	"github.com/nem0z/bitcoin-crawler/utils"
)

func (peer *Peer) Version() error {
	payload, err := payload.CreateVersion(peer.ip, peer.port)
	version, err := message.Create("version", payload.ToByte())
	if err != nil {
		return err
	}

	return peer.Send(version.ToByte())
}

func (peer *Peer) Verack() error {
	verack, err := message.Create("verack", []byte{})
	if err != nil {
		return err
	}

	return peer.Send(verack.ToByte())
}

func (peer *Peer) Ping(nonce []byte) (err error) {
	if nonce == nil {
		nonce, err = utils.CreateNonce(8)
		if err != nil {
			return err
		}
	}

	ping, err := message.Create("ping", nonce)
	if err != nil {
		return err
	}

	return peer.Send(ping.ToByte())
}

func (peer *Peer) GetAddr() error {
	ping, err := message.Create("getaddr", []byte{})
	if err != nil {
		return err
	}

	return peer.Send(ping.ToByte())
}

func (peer *Peer) Mempool() error {
	ping, err := message.Create("mempool", []byte{})
	if err != nil {
		return err
	}

	return peer.Send(ping.ToByte())
}
