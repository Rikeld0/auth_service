package connectorDB

import (
	"auth/pkg/config"
	"github.com/jackc/pgx"
	"sync"
)

type SQL interface {
	Exec(query string, args ...interface{}) error
	Close()
}

type Postgres struct {
	SQL
	sync.Mutex
	conn *pgx.Conn
}

func (db *Postgres) Exec(query string, args ...interface{}) error {
	db.Lock()
	_, err := db.conn.Exec(query, args)
	if err != nil {
		return err
	}
	db.Unlock()
	return nil
}

func ConnOpen() (*Postgres, error) {
	conn := &Postgres{}
	var err error
	connConfig, _ := pgx.ParseDSN(config.ConnInfo())
	conn.conn, err = pgx.Connect(connConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
