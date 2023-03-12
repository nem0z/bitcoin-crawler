package handlers

import (
	"encoding/binary"
	"net"

	"github.com/nem0z/bitcoin-crawler/peer"
)

func ParseAddr(addr []byte) *peer.Addr {
	ip := net.IP(addr[12:28]).To16()
	port := binary.BigEndian.Uint16(addr[28:30])

	return &peer.Addr{Ip: ip.String(), Port: int(port)}
}

func ParseListAddr(payload []byte) []*peer.Addr {
	_, n := binary.Uvarint(payload)
	addrs := make([]*peer.Addr, (len(payload)-n)/30)

	for i := range addrs {
		addrs[i] = ParseAddr(payload[i*30+n : 30*(i+1)+n])
	}

	return addrs
}
