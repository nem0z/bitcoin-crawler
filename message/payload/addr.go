package payload

import (
	"bytes"
	"encoding/binary"
	"net"
)

type Addr struct {
	Timestamp uint32
	Service   uint64
	Ip        string
	Port      uint16
}

func parseAddr(addr []byte) *Addr {
	timestamp := binary.LittleEndian.Uint32(addr[:4])
	services := binary.LittleEndian.Uint64(addr[4:12])
	ip := net.IP(addr[12:28]).To16()
	port := binary.BigEndian.Uint16(addr[28:30])

	return &Addr{timestamp, services, ip.String(), port}
}

func ParseAddr(payload []byte) []*Addr {
	_, n := binary.Uvarint(payload)
	addrs := make([]*Addr, (len(payload)-n)/30)

	for i := range addrs {
		if bytes.Equal(payload[i*30+n:30*(i+1)+n], make([]byte, 30)) {
			if i <= 1 {
				return []*Addr{}
			}
			return addrs[:i-1]
		}
		addrs[i] = parseAddr(payload[i*30+n : 30*(i+1)+n])
	}

	return addrs
}
