package db_test

import (
	"simplewallet/util/db"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestInitDb(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks
	t.Run("init db success", func(t *testing.T) {
		err := db.InitDbMock()
		assert.Nil(t, err)
		cli := db.GetDbClientMock()
		assert.NotNil(t, cli)
		cli.Close()
	})
}
