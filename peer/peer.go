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
	ip        string
	port      int
	conn      net.Conn
	handlers  Handlers
	Info      *Info
	Addrs     []*Addr
	PingNonce []byte
	PingAt    time.Time
	PongAt    time.Time
	queue     chan *message.Message
}

// Create the net.coon with the peer
func New(ip string, port int) (*Peer, error) {
	var err error
	var conn net.Conn

	netIp := net.ParseIP(ip)
	if netIp.To4() == nil {
		conn, err = net.DialTimeout("tcp6", fmt.Sprintf("[%v]:%v", ip, port), time.Second*3)
	} else {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Second*3)
	}

	ch := make(chan *message.Message)

	peer := &Peer{ip, port, conn, Handlers{}, nil, []*Addr{}, nil, time.Time{}, time.Time{}, ch}
	peer.start()

	return peer, err
}

func (peer *Peer) start() {
	go peer.Handle()
}

// Send a message to the peer
func (peer *Peer) Send(msg *message.Message) error {
	if !msg.IsValid() {
		return errors.New("Trying to send an invalid message")
	}

	msgData := msg.MarshalMessage()
	n, err := peer.conn.Write(msgData)

	if err != nil {
		return err
	}

	if n != len(msgData) {
		return errors.New("Wrong number of bytes sent")
	}

	log.Println("Sending message :", string(msg.Command))
	return nil
}

func (peer *Peer) Queue(msg *message.Message) {
	peer.queue <- msg
}

func (peer *Peer) ConsumeQueue() {

	go func() {
		for msg := range peer.queue {
			err := peer.Send(msg)
			if err != nil {
				log.Println("Consuming queue :", err)
			}
		}
	}()
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
	lengthValue := binary.LittleEndian.Uint32(length)
	checksum := header[20:24]

	payload := make([]byte, lengthValue)

	if lengthValue > 0 {

		totalRead := uint32(0)

		for totalRead < lengthValue {

			n, err := peer.conn.Read(payload[totalRead:])
			if err != nil {
				return nil, err
			}

			totalRead += uint32(n)
			if n == 0 {
				return nil, errors.New("Read message : No more byte to read")
			}
		}
	}

	return &message.Message{
		MagicNo:  magic,
		Command:  command,
		Length:   length,
		Checksum: checksum,
		Payload:  payload,
	}, err
}

func (peer *Peer) SelfAddr() *Addr {
	return &Addr{peer.ip, peer.port}
}
