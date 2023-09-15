package redisdb

import (
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
		Addr:     os.Getenv("ADR_REDIS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	return client
}
