package util_test

import (
	"context"
	"log"
	"simplewallet/util"
	"simplewallet/util/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func ConnectRedis() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := db.InitRedisMock()
	if err != nil {
		panic(err)
	}
}
func DisconnectRedis() {
	if err := db.GetRedisClientMock().Close(); err != nil {
		log.Println("fail to disconnect redis", err)
	}
}
func TestLock(t *testing.T) {
	ConnectRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectRedis()

	t.Run("testLock sussess", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		key := "testLock:" + logID
		locker := util.NewDistributedLockMock(logID, key, 10, ctx)
		err := locker.Lock()
		assert.Nil(t, err)

		rsp, err := db.GetRedisClientMock().Do(ctx, "EXISTS", key).Result()
		assert.Nil(t, err)
		assert.Equal(t, rsp, int64(1))

		locker.UnLock()
	})

	t.Run("testLock fail-[expireTime < 0]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		key := "testLock:" + logID
		locker := util.NewDistributedLockMock(logID, key, -1, ctx)
		err := locker.Lock()
		assert.NotNil(t, err)
	})

}

func TestUnlock(t *testing.T) {
	ConnectRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectRedis()

	t.Run("testUnLock sussess", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		key := "testUnLock:" + logID
		locker := util.NewDistributedLockMock(logID, key, 10, ctx)
		err := locker.Lock()
		assert.Nil(t, err)

		rsp, err := db.GetRedisClientMock().Do(ctx, "EXISTS", key).Result()
		assert.Nil(t, err)
		assert.Equal(t, rsp, int64(1))

		locker.UnLock()
		rsp, err = db.GetRedisClientMock().Do(ctx, "EXISTS", key).Result()
		assert.Nil(t, err)
		assert.Equal(t, rsp, int64(0))
	})

	t.Run("testUnLock success-[repeated unlock]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		key := "testUnLock:" + logID
		locker := util.NewDistributedLockMock(logID, key, 10, ctx)
		err := locker.Lock()
		assert.Nil(t, err)
		err = locker.UnLock()
		assert.Nil(t, err)
		err = locker.UnLock() // repeated unlock
		assert.Nil(t, err)
	})

}
