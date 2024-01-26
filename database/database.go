package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

const (
	createTableNodes = `
		CREATE TABLE IF NOT EXISTS nodes (
			id INTEGER PRIMARY KEY,
			ip VARCHAR(255),
			port INTEGER,
			version INTEGER,
			services INTEGER,
			relay BOOLEAN
		)
	`

	createTablePings = `
		CREATE TABLE IF NOT EXISTS pings (
			id INTEGER PRIMARY KEY,
			node_id INTEGER,
			timestamp TIMESTAMP,
			ok BOOLEAN,
			FOREIGN KEY (node_id) REFERENCES nodes(id)
		)
	`
)

func Init(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableNodes)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTablePings)
	return &DB{db}, nil
}
