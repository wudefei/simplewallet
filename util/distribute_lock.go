package util

import (
	"errors"
	"log"
	"simplewallet/util/db"
)

type DistributedLock interface {
	Lock() error
	UnLock() error
}

type RedisDistributedLock struct {
	LogId    string
	Key      string
	Expire   int
	isLocked bool
}

func NewDistributedLock(logId string, key string, expire int) DistributedLock {
	return &RedisDistributedLock{LogId: logId, Key: key, Expire: expire, isLocked: false}
}

func (r *RedisDistributedLock) Lock() error {
	cli := db.GetRedisClient()
	resp, err := cli.Do("SET", r.Key, "", "NX", "EX", r.Expire).Result()
	if err != nil {
		return err
	}
	log.Printf("%s|get lock success:%s\n", r.LogId, resp)
	if resp != "OK" {
		return errors.New("lock key fail")
	}
	r.isLocked = true
	return nil
}

func (r *RedisDistributedLock) UnLock() error {
	cli := db.GetRedisClient()
	if !r.isLocked {
		return nil
	}
	resp, err := cli.Do("DEL", r.Key).Result()
	if err != nil {
		return err
	}
	r.isLocked = false
	log.Printf("%s|unlock key success:%s\n", r.LogId, resp)
	return nil
}
