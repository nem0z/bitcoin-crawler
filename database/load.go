package database

import (
	"github.com/nem0z/bitcoin-crawler/peer"
)

func (db *DB) LoadAddrs(ok bool) ([]*peer.Addr, error) {
	req := `SELECT n.ip, n.port
			FROM nodes n
			JOIN pings p ON p.node_id = n.id
			JOIN (
				SELECT node_id, MAX(timestamp) AS latest_timestamp
				FROM pings
				WHERE ok = true
				GROUP BY node_id
			) latest_ping ON p.node_id = latest_ping.node_id AND p.timestamp = latest_ping.latest_timestamp
			WHERE p.ok = true;`

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
