package db

import (
	"github.com/go-redis/redis/v8"
)

var redisClientMock *redis.Client

func InitRedisMock() error {
	redisClientMock = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
	return nil
}

func GetRedisClientMock() *redis.Client {
	return redisClientMock
}
