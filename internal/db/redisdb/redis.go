package redisdb

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RDB struct {
	client *redis.Client
}

func NewClientRedis(client *redis.Client) RedisClient {
	return &RDB{client: client}
}

func (c *RDB) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(ctx, key)
}

func (c *RDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}

func (c *RDB) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.client.Del(ctx, keys...)
}

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
