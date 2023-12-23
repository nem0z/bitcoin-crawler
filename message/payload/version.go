package payload

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
	Version     [4]byte
	Services    [8]byte
	Timestamp   [8]byte
	AddrPeer    []byte
	AddrClient  []byte
	Nonce       [8]byte
	UserAgent   []byte
	StartHeight [4]byte
	Relay       byte
}

func NewVersion(ip string, port int) (*Version, error) {
	msgVersion := new(Version)

	binary.LittleEndian.PutUint32(msgVersion.Version[:], defaultProtocolVersion)
	binary.LittleEndian.PutUint64(msgVersion.Services[:], 1)
	binary.LittleEndian.PutUint64(msgVersion.Timestamp[:], uint64(time.Now().Unix()))
	binary.LittleEndian.PutUint64(msgVersion.Nonce[:], 0)
	binary.LittleEndian.PutUint32(msgVersion.StartHeight[:], defaultStartHeight)

	clientIp := net.IPv6loopback.To16()
	clientPort := make([]byte, 2)
	binary.BigEndian.PutUint16(clientPort, defaultClientPort)

	peerIp := net.ParseIP(ip).To16()
	peerPort := make([]byte, 2)
	binary.BigEndian.PutUint16(peerPort, uint16(port))

	addrClient := bytes.Join([][]byte{
		msgVersion.Services[:],
		clientIp,
		clientPort,
	}, []byte{})

	addrPeer := bytes.Join([][]byte{
		msgVersion.Services[:],
		peerIp,
		peerPort,
	}, []byte{})

	msgVersion.AddrClient = addrClient
	msgVersion.AddrPeer = addrPeer
	msgVersion.Relay = 0

	return msgVersion, nil
}

func (v *Version) ToByte() []byte {
	return bytes.Join([][]byte{
		v.Version[:],
		v.Services[:],
		v.Timestamp[:],
		v.AddrClient,
		v.AddrPeer,
		v.Nonce[:],
		v.UserAgent,
		v.StartHeight[:],
		{v.Relay},
	}, []byte{})
}
