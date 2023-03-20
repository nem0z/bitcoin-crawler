package peer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nem0z/bitcoin-crawler/message"
)

type Info struct {
	Version  int32
	Services uint64
	Relay    bool
}

type Addr struct {
	Ip   string
	Port int
}

type Peer struct {
	ip       string
	port     int
	conn     net.Conn
	stop     chan bool
	handlers Handlers
	Info     *Info
}

// Create the net.coon with the peer
func Create(ip string, port int, stop chan bool) (*Peer, error) {
	var err error
	var conn net.Conn
	timeout := time.Second * 5

	netIp := net.ParseIP(ip)
	if netIp.To4() == nil {
		conn, err = net.DialTimeout("tcp6", fmt.Sprintf("[%v]:%v", ip, port), timeout)
	} else {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), timeout)
	}

	return &Peer{ip, port, conn, stop, Handlers{}, &Info{}}, err
}

// Init the connection with the peer by sending version and verack message
func (peer *Peer) Init() error {
	go peer.Handle()

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
	peer.stop <- false
	peer.conn.Close()
	peer = nil
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

// Return Addr of himself (as ip and port are private)
func (peer *Peer) SelfAddr() *Addr {
	return &Addr{peer.ip, peer.port}
}
