package redis

import (
	"context"
	"fmt"
	"niko-web_app/settings"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password, // 密码
		DB:       cfg.DB,       // 数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})
	_, err := rdb.Ping(context.Background()).Result()
	return err
}

func Close() {
	_ = rdb.Close()
}
