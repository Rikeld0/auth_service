package connector_db

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4"
	"time"
)

type Postgre interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) error
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Close()
}

type Redis interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(key string) error
	Close() error
}
