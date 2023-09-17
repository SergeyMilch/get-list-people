package redisdb

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("ADDR_REDIS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	// Проверка подключения к Redis
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}
