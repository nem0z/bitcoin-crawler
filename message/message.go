package message

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/nem0z/bitcoin-crawler/utils"
)

var mainMagicNo = []byte{0xf9, 0xbe, 0xb4, 0xD9}

var availableCommands = []string{
	"version",
	"verack",
	"ping",
	"pong",
	"getaddr",
	"mempool",
}

type Message struct {
	MagicNo  []byte
	Command  []byte
	Length   []byte
	Checksum []byte
	Payload  []byte
}

func New(commandName string, payload []byte) (*Message, error) {
	command, err := formatCommandName(commandName)
	if err != nil {
		return nil, err
	}

	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(payload)))

	checksum := utils.Checksum(payload)

	return &Message{
		mainMagicNo,
		command,
		length,
		checksum,
		payload,
	}, nil
}

func (m *Message) IsValid() bool {
	checksum := utils.Checksum(m.Payload)
	length := binary.LittleEndian.Uint32(m.Length)

	return bytes.Equal(m.MagicNo, mainMagicNo) &&
		uint32(len(m.Payload)) == length &&
		bytes.Equal(m.Checksum, checksum)
}

}

func (m *Message) Display() {
	lenghtValue := binary.LittleEndian.Uint32(m.Length)
	checksum := utils.Checksum(m.Payload)
	valid := fmt.Sprintf("%x", checksum) == fmt.Sprintf("%x", m.Checksum)

	fmt.Println("*----------*")
	fmt.Println("Magic number :", m.MagicNo)
	fmt.Printf("Command : %s\n", m.Command)
	fmt.Println("Length :", lenghtValue)
	fmt.Printf("Checksum (%v) : %v\n", valid, m.Checksum)
	fmt.Println("Payload length :", len(m.Payload))
	fmt.Printf("*----------*\n\n")
}

func (m *Message) ToByte() []byte {
	return bytes.Join([][]byte{
		m.MagicNo,
		m.Command,
		m.Length,
		m.Checksum,
		m.Payload,
	}, []byte{})
}
