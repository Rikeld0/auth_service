package connector_db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type connPosgres struct {
	conn *pgxpool.Pool
}

func NewPostreConn(conn *pgxpool.Pool) Postgre {
	c := &connPosgres{
		conn: conn,
	}
	return c
}

func (c *connPosgres) Exec(ctx context.Context, query string, args ...interface{}) error {
	if _, err := c.conn.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (c *connPosgres) Query(ctx context.Context, query string, args ...interface{}) error {
	panic("implement me")
}

func (c *connPosgres) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return c.conn.QueryRow(ctx, query, args...)
}

func (c *connPosgres) Close() {
	c.conn.Close()
}

func ConnPostger(config string) (Postgre, error) {
	connConfig, _ := pgxpool.ParseConfig(config)
	conn, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("error:", err)
		return nil, err
	}
	return NewPostreConn(conn), nil
}
