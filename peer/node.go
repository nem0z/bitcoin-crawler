package peer

import "time"

type Node struct {
	Timestamp time.Time
	Info      *Info
	Addr      *Addr
	Ping      bool
}
