package client

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"watchAlert/public/globals"
)

func InitRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", globals.Config.Redis.Host, globals.Config.Redis.Port),
		Password: globals.Config.Redis.Pass,
		DB:       0, // 使用默认的数据库
	})

	// 尝试连接到 Redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("redis Connection Failed %s", err)
		return nil
	}

	return client

}
