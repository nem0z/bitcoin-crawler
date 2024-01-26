package database

import (
	"github.com/nem0z/bitcoin-crawler/peer"
)

func (db *DB) LoadAddrs(ok bool) ([]*peer.Addr, error) {
	req := `SELECT n.ip, n.port 
			FROM nodes n 
			JOIN pings p
			ON p.node_id = n.id
			WHERE p.ok = ?`

	rows, err := db.Query(req, ok)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addrs := []*peer.Addr{}
	for rows.Next() {
		addr := &peer.Addr{}
		err := rows.Scan(&addr.Ip, &addr.Port)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}

	if err := rows.Err(); err != nil {
		return addrs, err
	}

	return addrs, nil
}
