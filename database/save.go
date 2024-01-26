package database

import (
	"database/sql"
	"time"

	"github.com/nem0z/bitcoin-crawler/peer"
)

func (db *DB) GetNodeID(node *peer.Node) (id int64, err error) {
	req := "SELECT id FROM nodes WHERE ip = ? AND port = ?"
	err = db.QueryRow(req, node.Addr.Ip, node.Addr.Port).Scan(&id)
	return id, err
}

func (db *DB) createNode(node *peer.Node) (id int64, err error) {
	var result sql.Result
	req := `INSERT INTO nodes (ip, port, version, services, relay)
			VALUES (?, ?, ?, ?, ?)`

	if node.Info != nil {
		result, _ = db.Exec(req, node.Addr.Ip, node.Addr.Port, node.Info.Version, node.Info.Services, node.Info.Relay)
	} else {
		result, _ = db.Exec(req, node.Addr.Ip, node.Addr.Port, nil, nil, nil)
	}

	return result.LastInsertId()
}

func (db *DB) updateNode(id int64, node *peer.Node) error {
	req := `UPDATE nodes
			SET version = ?, services = ?, relay = ?
			WHERE id = ?`

	if node.Info == nil {
		return nil
	}

	_, err := db.Exec(req, node.Info.Version, node.Info.Services, node.Info.Relay, id)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) CreateOrUpdateNode(node *peer.Node) (int64, error) {
	id, err := db.GetNodeID(node)

	if err == sql.ErrNoRows {
		return db.createNode(node)
	}

	return id, db.updateNode(id, node)
}

func (db *DB) CreatePing(id int64, ok bool) error {
	req := `INSERT INTO pings (node_id, ok, timestamp)
	VALUES (?, ?, ?)`

	_, err := db.Exec(req, id, ok, time.Now())
	return err
}

func (db *DB) Update(node *peer.Node) error {
	id, err := db.CreateOrUpdateNode(node)
	if err != nil {
		return err
	}

	return db.CreatePing(id, node.Ping)
}
