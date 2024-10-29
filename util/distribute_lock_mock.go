package util

import (
	"context"
	"errors"
	"log"
	"simplewallet/util/db"
	"time"
)

type RedisDistributedLockMock struct {
	LogId    string
	Key      string
	Expire   int
	isLocked bool
	ctx      context.Context
}

func NewDistributedLockMock(logId string, key string, expire int, ctx context.Context) DistributedLock {
	return &RedisDistributedLockMock{LogId: logId, Key: key, Expire: expire, isLocked: false, ctx: ctx}
}

func (r *RedisDistributedLockMock) Lock() error {
	cli := db.GetRedisClientMock()
	expire := time.Second * time.Duration(r.Expire)
	_, err := cli.SetNX(r.ctx, r.Key, "", expire).Result()
	if err != nil {
		return err
	}
	r.isLocked = true
	log.Printf("%s|get lock success\n", r.LogId)
	return nil
}

func (r *RedisDistributedLockMock) UnLock() error {
	cli := db.GetRedisClientMock()
	if !r.isLocked {
		return nil
	}
	resp, err := cli.Del(r.ctx, r.Key).Result()
	if err != nil {
		return err
	}
	if resp == 0 {
		return errors.New("unlock key fail")
	}
	r.isLocked = false
	log.Printf("%s|unlock key success:%d\n", r.LogId, resp)
	return nil
}
