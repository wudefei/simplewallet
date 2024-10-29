package db

import (
	"errors"

	"github.com/go-redis/redis"
)

type RedisConf struct {
	Uri      string `yaml:"uri" json:"uri"`
	Password string `yaml:"password" json:"password"`
	Db       int    `yaml:"db" json:"db"`
}

var redisClient *redis.Client

func InitRedis(conf *RedisConf) error {
	if conf == nil {
		return errors.New("redis config is nil")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Uri,
		Password: conf.Password,
		DB:       conf.Db,
	})
	return nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}
