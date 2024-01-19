package peer

import (
	"time"
)

type Node struct {
	Timestamp time.Time `json:"timestamp"`
	Info      *Info     `json:"info"`
	Addr      *Addr     `json:"addr"`
	Ping      bool      `json:"ping"`
}
