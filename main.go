package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/nem0z/bitcoin-crawler/message"
	"github.com/nem0z/bitcoin-crawler/message/messages"
	"github.com/nem0z/bitcoin-crawler/utils"
)

func Read(conn net.Conn) (*message.Message, error) {
	header := make([]byte, 24)
	_, err := conn.Read(header)
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
		_, err = conn.Read(payload)
		if err != nil {
			return nil, err
		}
	}

	return &message.Message{
		MagicNo:  magic,
		Command:  command,
		Length:   length,
		Checksum: checksum,
		Payload:  payload,
	}, nil

}

func main() {
	peerIp := "2a02:8108:8ac0:207b:d250:99ff:fe9e:792a"
	peerPort := 8333
	timeout := 5 * time.Second

	versionPayload, err := messages.CreateVersion(peerIp, peerPort)
	utils.Handle(err)

	versionMessage, err := message.Create("version", versionPayload.ToByte())
	utils.Handle(err)

	conn, err := net.DialTimeout("tcp6", fmt.Sprintf("[%v]:%v", peerIp, peerPort), timeout)

	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Version")
	conn.Write(versionMessage.ToByte())
	msg, err := Read(conn)
	utils.Handle(err)
	msg.Display()

	fmt.Println("Verack")
	verackMsg, err := message.Create("verack", []byte{})
	utils.Handle(err)
	conn.Write(verackMsg.ToByte())
	msg, err = Read(conn)
	utils.Handle(err)
	msg.Display()

	fmt.Println("Ping")
	data := make([]byte, 8)
	pingMsg, err := message.Create("ping", data)
	utils.Handle(err)
	conn.Write(pingMsg.ToByte())

	fmt.Println("Get addr")
	getaddrMsg, err := message.Create("getaddr", []byte{})
	utils.Handle(err)
	conn.Write(getaddrMsg.ToByte())

	fmt.Println("Mempool")
	mempoolMsg, err := message.Create("mempool", []byte{})
	utils.Handle(err)
	conn.Write(mempoolMsg.ToByte())

	for i := 0; i < 10; i++ {
		msg, err = Read(conn)
		if err != nil {
			break
		}
		msg.Display()
	}

}
