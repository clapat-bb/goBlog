package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(host string, port string, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})
	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}
	log.Println("redis 连接成功")
}
