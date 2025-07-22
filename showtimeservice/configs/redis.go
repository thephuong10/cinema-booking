package configs

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}

	return rdb
}
