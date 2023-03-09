package peer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nem0z/bitcoin-crawler/message"
)

const defaultTime = 5 * time.Second

type Peer struct {
	ip   string
	port int
	conn net.Conn
}

// Create the net.coon with the peer
func Create(ip string, port int) (*Peer, error) {
	conn, err := net.DialTimeout("tcp6", fmt.Sprintf("[%v]:%v", ip, port), defaultTime)
	return &Peer{ip, port, conn}, err
}

// Init the connection with the peer by sending version and verack message
func (peer *Peer) Init() error {

	if err := peer.Version(); err != nil {
		return err
	}

	if err := peer.Verack(); err != nil {
		return err
	}

	return nil
}

// Close the conn with the peer
func (peer *Peer) Close() {
	peer.conn.Close()
}

// Send a message to the peer
func (peer *Peer) Send(message []byte) error {
	n, err := peer.conn.Write(message)

	if err != nil {
		return err
	}

	if n != len(message) {
		return errors.New("Wrong number of bytes sent")
	}

	return nil
}

// Read a message from the conn
func (peer *Peer) Read() (*message.Message, error) {
	header := make([]byte, 24)
	_, err := peer.conn.Read(header)
	if err != nil {
		return nil, err
	}

	magic := header[:4]
	command := header[4:16]
	length := header[16:20]
	lenghtValue := binary.LittleEndian.Uint32(length)
	checksum := header[20:24]

	payload := make([]byte, lenghtValue)

	if lenghtValue > 0 {
		_, err = peer.conn.Read(payload)
	}

	return &message.Message{
		MagicNo:  magic,
		Command:  command,
		Length:   length,
		Checksum: checksum,
		Payload:  payload,
	}, err
}
