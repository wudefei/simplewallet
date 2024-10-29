package db_test

import (
	"simplewallet/util/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestInitRedis(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks
	t.Run("init redis sussess", func(t *testing.T) {
		err := db.InitRedisMock()
		assert.Nil(t, err)
		redisCli := db.GetRedisClientMock()
		assert.NotNil(t, redisCli)
		redisCli.Close()
	})
}
