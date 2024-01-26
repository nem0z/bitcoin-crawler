package peer

import (
	"time"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/message/payload"
)

func (peer *Peer) Version() error {
	payload, err := payload.NewVersion(peer.ip, peer.port)
	msg, err := message.New("version", payload.ToByte())
	if err != nil {
		return err
	}

	return peer.Send(msg)
}

func (peer *Peer) Verack() error {
	msg, err := message.New("verack", []byte{})
	if err != nil {
		return err
	}

	return peer.Send(msg)
}

func (peer *Peer) Ping() error {
	nonce, err := message.CreateNonce(8)
	if err != nil {
		return err
	}

	msg, err := message.New("ping", nonce)
	if err != nil {
		return err
	}

	peer.Queue(msg)

	peer.PingAt = time.Now()
	peer.PingNonce = nonce

	return nil
}

func (peer *Peer) Pong(nonce []byte) error {
	msg, err := message.New("pong", nonce)
	if err != nil {
		return err
	}

	peer.Queue(msg)
	return nil
}

func (peer *Peer) GetAddr() error {
	msg, err := message.New("getaddr", []byte{})
	if err != nil {
		return err
	}

	peer.Queue(msg)
	return nil
}

func (peer *Peer) GetAddrNoError() {
	msg, err := message.New("getaddr", []byte{})
	if err != nil {
		return
	}

	peer.Queue(msg)
}

func (peer *Peer) Mempool() error {
	msg, err := message.New("mempool", []byte{})
	if err != nil {
		return err
	}

	peer.Queue(msg)
	return nil
}
