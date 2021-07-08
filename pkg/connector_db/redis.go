package connector_db

import (
	"github.com/go-redis/redis"
	"time"
)

type myRedis struct {
	rClient *redis.Client
}

func newMyRedis(rClient *redis.Client) Redis {
	return &myRedis{
		rClient: rClient,
	}
}

func (m *myRedis) Get(key string) *redis.StringCmd {
	return m.rClient.Get(key)
}

func (m *myRedis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.rClient.Set(key, value, expiration)
}

func (m *myRedis) Close() error {
	return m.rClient.Close()
}

func ConnRedis(opt *redis.Options) Redis {
	conn := redis.NewClient(opt)
	return newMyRedis(conn)
}
