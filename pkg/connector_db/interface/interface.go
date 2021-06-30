package _interface

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) error
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Close()
}
