package messages

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

const defaultProtocolVersion = 70016
const defaultStartHeight = 0
const defaultClientPort = 8333

type Version struct {
	Version     []byte
	Services    []byte
	Timestamp   []byte
	AddrPeer    []byte
	AddrClient  []byte
	Nonce       []byte
	UserAgent   []byte
	StartHeight []byte
	Relay       []byte
}

func CreateVersion(ip string, port int) (*Version, error) {
	version := make([]byte, 4)
	services := make([]byte, 8)
	timestamp := make([]byte, 8)
	nonce := make([]byte, 8)
	startHeight := make([]byte, 4)

	binary.LittleEndian.PutUint32(version, 70016)
	binary.LittleEndian.PutUint32(services, 1)
	binary.LittleEndian.PutUint32(timestamp, uint32(time.Now().Unix()))
	binary.LittleEndian.PutUint32(nonce, 0)
	binary.LittleEndian.PutUint32(startHeight, defaultStartHeight)

	clientIp := net.IPv6loopback.To16()
	clientPort := make([]byte, 2)
	binary.BigEndian.PutUint16(clientPort, defaultClientPort)

	peerIp := net.ParseIP(ip).To16()
	peerPort := make([]byte, 2)
	binary.BigEndian.PutUint16(peerPort, uint16(port))

	addrClient := bytes.Join([][]byte{
		services,
		clientIp,
		clientPort,
	}, []byte{})

	addrPeer := bytes.Join([][]byte{
		services,
		peerIp,
		peerPort,
	}, []byte{})

	return &Version{
		version,
		services,
		timestamp,
		addrClient,
		addrPeer,
		nonce,
		[]byte{0},
		startHeight,
		[]byte{0},
	}, nil
}

func (v *Version) ToByte() []byte {
	return bytes.Join([][]byte{
		v.Version,
		v.Services,
		v.Timestamp,
		v.AddrClient,
		v.AddrPeer,
		v.Nonce,
		v.UserAgent,
		v.StartHeight,
		v.Relay,
	}, []byte{})
}
