package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("127.0.0.1:6379"),
		// Password: "",               // No password set
		DB: 0, // Use default DB
	})
}
func GetRedis() *redis.Client {
	return redisClient

}
