package peer

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nem0z/bitcoin-crawler/message"
)

type Info struct {
	Version  int32   `json:"version"`
	Services uint64  `json:"services"`
	Relay    bool    `json:"relay"`
	Addrs    []*Addr `json:"addrs"`
}

type PingInfo struct {
	PingNonce []byte
	PingAt    time.Time
	PongAt    time.Time
}

type Addr struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func (addr *Addr) String() string {
	return fmt.Sprintf("%v:%v", addr.Ip, addr.Port)
}

type Peer struct {
	addr     *Addr
	conn     net.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	handlers Handlers
	queue    chan *message.Message
	onClose  chan *Node
	Info     *Info
	PingInfo *PingInfo
}

func (p *Peer) Addr() string {
	return p.addr.String()
}

// Done will block the goroutine calling this function until the peer ctx timeout or get canceled
func (p *Peer) Done() {
	<-p.ctx.Done()
}

// Create the net.coon with the peer
func New(addr *Addr, onClose chan *Node) (*Peer, error) {
	var err error
	var conn net.Conn

	netIp := net.ParseIP(addr.Ip)
	if netIp.To4() == nil {
		conn, err = net.DialTimeout("tcp6", fmt.Sprintf("[%v]:%v", addr.Ip, addr.Port), time.Second*3)
	} else {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%v:%v", addr.Ip, addr.Port), time.Second*3)
	}

	if err != nil {
		return nil, err
	}

	queue := make(chan *message.Message, 100)

	return &Peer{
		addr: addr,
		conn: conn, handlers: Handlers{},
		ctx:      context.Background(),
		queue:    queue,
		onClose:  onClose,
		Info:     &Info{},
		PingInfo: &PingInfo{},
	}, nil
}

func (peer *Peer) Start() error {
	ctx, cancel := context.WithTimeout(peer.ctx, time.Second*30)
	peer.ctx = ctx
	peer.cancel = cancel

	go func() {
		<-ctx.Done()
		peer.Close()
	}()

	if err := peer.Version(); err != nil {
		return err
	}

	if err := peer.Verack(); err != nil {
		return err
	}

	if err := peer.Ping(); err != nil {
		return err
	}

	if err := peer.GetAddr(); err != nil {
		return err
	}

	go peer.Handle()
	return nil
}

// Send a message to the peer
func (peer *Peer) Send(msg *message.Message) error {
	if !msg.IsValid() {
		return errors.New("trying to send an invalid message")
	}

	msgData := msg.MarshalMessage()
	n, err := peer.conn.Write(msgData)

	if err != nil {
		return err
	}

	if n != len(msgData) {
		return errors.New("wrong number of bytes sent")
	}

	// log.Println("Send message :", string(msg.Command))
	return nil
}

func (peer *Peer) Queue(msg *message.Message) {
	peer.queue <- msg
}

func (peer *Peer) ConsumeQueue() {
	for msg := range peer.queue {
		select {
		case <-peer.ctx.Done():
			return
		default:
			err := peer.Send(msg)
			if err != nil {
				// log.Println("Consuming queue :", err)
			}
		}
	}
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

func (peer *Peer) Display() {
	fmt.Println("/-----*-----/")
	fmt.Printf("Peer : %v:%v\n", peer.addr.Ip, peer.addr.Port)
	fmt.Printf("Info : %v : %v : %v\n", peer.Info.Version, peer.Info.Services, peer.Info.Relay)
	fmt.Println("Addrs : ", len(peer.Info.Addrs))
	fmt.Printf("Ping/Pong : %v => %v\n", peer.PingInfo.PingAt, peer.PingInfo.PongAt)
	fmt.Printf("/-----*-----/\n\n")
}

func (peer *Peer) Close() {
	peer.cancel()
	if peer.conn == nil {
		return
	}

	node := &Node{
		time.Now(),
		peer.Info,
		peer.addr,
		peer.PingInfo.PongAt.Sub(peer.PingInfo.PingAt) > 0,
	}

	peer.onClose <- node

	peer.conn.Close()
}
