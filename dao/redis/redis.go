package redis

import (
	"bluebell/settings"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

type SliceCmd = redis.SliceCmd
type StringStringMapCmd = redis.StringStringMapCmd

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return
	}
	return
}

func Close() {
	_ = client.Close()
}
