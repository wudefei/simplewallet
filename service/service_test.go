package service_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"simplewallet/data"
	"simplewallet/service"
	"simplewallet/util"
	"simplewallet/util/db"
	"simplewallet/util/errcode"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func ConnectDBRedis() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := db.InitDbMock()
	if err != nil {
		panic(err)
	}
	err = db.InitRedisMock()
	if err != nil {
		panic(err)
	}
}
func DisconnectDBRedis() {
	dbCli := db.GetDbClientMock()
	redisCli := db.GetRedisClientMock()
	err := dbCli.Close()
	if err != nil {
		log.Println("fail to close db", err)
	}
	err = redisCli.Close()
	if err != nil {
		log.Println("fail to close redis", err)
	}
}

func TestDeposit(t *testing.T) {
	ConnectDBRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectDBRedis()

	mockDBCli := db.GetDbClientMock()
	mock := db.GetSqlMock()
	t.Run("case1: deposit success-[user wallet not exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "deposit:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)

		depositReq := &data.DepositReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(depositReq.OrderID).WillReturnRows(transRows)

		walletRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(depositReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(depositReq.UserID, depositReq.Amount, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(depositReq.OrderID, depositReq.UserID, data.TxTypeDeposit, depositReq.Amount, 0, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Deposit(depositReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case2: deposit success-[user wallet exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "deposit:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)

		depositReq := &data.DepositReq{OrderID: logID, UserID: 101, Amount: 2000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(depositReq.OrderID).WillReturnRows(transRows)

		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, depositReq.UserID, 1000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(depositReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance + $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(depositReq.Amount, tn, depositReq.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(depositReq.OrderID, depositReq.UserID, data.TxTypeDeposit, depositReq.Amount, 0, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Deposit(depositReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case3: deposit fail-[insert wallet fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "deposit:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)

		depositReq := &data.DepositReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(depositReq.OrderID).WillReturnRows(transRows)

		walletRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(depositReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(depositReq.UserID, depositReq.Amount, tn, tn).WillReturnError(errors.New("insert wallet fail"))

		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Deposit(depositReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("case4: deposit fail-[insert transaction fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "deposit:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)

		depositReq := &data.DepositReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(depositReq.OrderID).WillReturnRows(transRows)

		walletRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(depositReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(depositReq.UserID, depositReq.Amount, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(depositReq.OrderID, depositReq.UserID, data.TxTypeDeposit, depositReq.Amount, 0, tn, tn).WillReturnError(errors.New("insert transaction fail"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Deposit(depositReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case5: deposit fail-[order_id exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "deposit:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)

		depositReq := &data.DepositReq{OrderID: logID, UserID: 101, Amount: 2000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{"id", "order_id", "user_id", "tx_type", "amount", "related_user_id", "created_at", "updated_at"})
		transRows.AddRow(1, depositReq.OrderID, depositReq.UserID, data.TxTypeDeposit, depositReq.Amount, 0, tn, tn)
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(depositReq.OrderID).WillReturnRows(transRows)
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Deposit(depositReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeOrderIDRepeat, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case6: deposit fail-[distribute lock fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "deposit:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, -1, ctx) // expire time < 0
		depositReq := &data.DepositReq{OrderID: logID, UserID: 101, Amount: 2000.00}
		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Deposit(depositReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeLockFail, rsp.Code)
	})

}

func TestWithdraw(t *testing.T) {
	ConnectDBRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectDBRedis()

	mockDBCli := db.GetDbClientMock()
	mock := db.GetSqlMock()
	t.Run("case1: withdraw success-[balance > amount]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, withdrawReq.UserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(withdrawReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(withdrawReq.Amount, tn, withdrawReq.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(withdrawReq.OrderID, withdrawReq.UserID, data.TxTypeWithdraw, withdrawReq.Amount, 0, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case2: withdraw success-[balance == amount]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, withdrawReq.UserID, 1000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(withdrawReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(withdrawReq.Amount, tn, withdrawReq.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(withdrawReq.OrderID, withdrawReq.UserID, data.TxTypeWithdraw, withdrawReq.Amount, 0, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case3: withdraw fail-[balance < amount]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, withdrawReq.UserID, 500.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(withdrawReq.UserID).WillReturnRows(walletRows)
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeBalanceNotEnough, rsp.Code) // balance not enough
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case4: withdraw fail-[user wallet not exists]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(withdrawReq.UserID).WillReturnRows(walletRows)
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeUserWalletNotExist, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case5: withdraw fail-[order_id already exists]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{"id", "order_id", "user_id", "tx_type", "amount", "related_user_id", "created_at", "updated_at"})
		transRows.AddRow(1, withdrawReq.OrderID, withdrawReq.UserID, data.TxTypeWithdraw, withdrawReq.Amount, 0, tn, tn)
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeOrderIDRepeat, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case6: withdraw fail-[update wallet fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, withdrawReq.UserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(withdrawReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(withdrawReq.Amount, tn, withdrawReq.UserID).WillReturnError(errors.New("update wallet fail"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case7: withdraw fail-[insert transaction fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(withdrawReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, withdrawReq.UserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(withdrawReq.UserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(withdrawReq.Amount, tn, withdrawReq.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(withdrawReq.OrderID, withdrawReq.UserID, data.TxTypeWithdraw, withdrawReq.Amount, 0, tn, tn).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case8: withdraw fail-[distributed lock fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "withdraw:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, -1, ctx)
		withdrawReq := &data.WithdrawReq{OrderID: logID, UserID: 101, Amount: 1000.00}
		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Withdraw(withdrawReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeLockFail, rsp.Code)
	})
}

func TestTransfer(t *testing.T) {
	ConnectDBRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectDBRedis()

	mockDBCli := db.GetDbClientMock()
	mock := db.GetSqlMock()
	t.Run("case1: transfer success-[balance > amount]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.FromUserID).WillReturnResult(sqlmock.NewResult(1, 1))

		walletRows2 := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.ToUserID).WillReturnRows(walletRows2)
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(transferReq.ToUserID, transferReq.Amount, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(transferReq.OrderID, transferReq.FromUserID, data.TxTypeTransferOut, transferReq.Amount, transferReq.ToUserID, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(transferReq.OrderID, transferReq.ToUserID, data.TxTypeTransferIn, transferReq.Amount, transferReq.FromUserID, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case2: transfer success-[balance == amount]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 1000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.FromUserID).WillReturnResult(sqlmock.NewResult(1, 1))

		walletRows2 := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.ToUserID).WillReturnRows(walletRows2)
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(transferReq.ToUserID, transferReq.Amount, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(transferReq.OrderID, transferReq.FromUserID, data.TxTypeTransferOut, transferReq.Amount, transferReq.ToUserID, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(transferReq.OrderID, transferReq.ToUserID, data.TxTypeTransferIn, transferReq.Amount, transferReq.FromUserID, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case3: transfer fail-[balance < amount]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 500.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeBalanceNotEnough, rsp.Code) // user balance not enough
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case4: transfer fail-[user wallet not exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}

		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeUserWalletNotExist, rsp.Code) // user wallet not exist
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case5: transfer fail-[update sender balance fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.FromUserID).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code) // update wallet fail
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case6: transfer fail-[update receiver balance fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.FromUserID).WillReturnResult(sqlmock.NewResult(1, 1))

		walletRowsRecv := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.ToUserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.ToUserID).WillReturnRows(walletRowsRecv)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance + $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.ToUserID).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code) // update wallet fail
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case7: transfer fail-[insert receiver wallet record fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.FromUserID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.ToUserID).WillReturnError(sql.ErrNoRows)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(transferReq.ToUserID, transferReq.Amount, tn, tn).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code) // update wallet fail
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case8: transfer fail-[insert transaction fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, 5, ctx)
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		tn := time.Now().Unix()
		// mock DB data
		mock.ExpectBegin()
		transRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = ?").WithArgs(transferReq.OrderID).WillReturnRows(transRows)
		walletRows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, transferReq.FromUserID, 2000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.FromUserID).WillReturnRows(walletRows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3")).
			WithArgs(transferReq.Amount, tn, transferReq.FromUserID).WillReturnResult(sqlmock.NewResult(1, 1))

		walletRows2 := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(transferReq.ToUserID).WillReturnRows(walletRows2)
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(transferReq.ToUserID, transferReq.Amount, tn, tn).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
			WithArgs(transferReq.OrderID, transferReq.FromUserID, data.TxTypeTransferOut, transferReq.Amount, transferReq.ToUserID, tn, tn).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeDbError, rsp.Code)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("case9: transfer fail-[distribute lock fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		lockKey := "transfer:" + logID
		loker := util.NewDistributedLockMock(logID, lockKey, -1, ctx) // lock time < 0
		transferReq := &data.TransferReq{OrderID: logID, FromUserID: 101, ToUserID: 102, Amount: 1000.00}
		walletService := service.NewWalletService(ctx, logID, mockDBCli, loker)
		rsp, err := walletService.Transfer(transferReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeLockFail, rsp.Code)
	})
}

func TestGetBalance(t *testing.T) {
	ConnectDBRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectDBRedis()

	mockDBCli := db.GetDbClientMock()
	mock := db.GetSqlMock()
	t.Run("case1: get balance success-[record exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		getBalanceReq := &data.GetBalanceReq{UserID: 101}
		tn := time.Now().Unix()
		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, getBalanceReq.UserID, 1000.00, tn, tn)
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ? ").WithArgs(getBalanceReq.UserID).WillReturnRows(rows)
		walletService := service.NewWalletService(ctx, logID, mockDBCli, nil)
		rsp, err := walletService.GetBalance(getBalanceReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		assert.Equal(t, 1000.00, rsp.Data.Balance)
	})

	t.Run("case2: get balance fail-[user wallet record not exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		getBalanceReq := &data.GetBalanceReq{UserID: 101}
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ?").WithArgs(getBalanceReq.UserID).WillReturnError(sql.ErrNoRows)
		walletService := service.NewWalletService(ctx, logID, mockDBCli, nil)
		rsp, err := walletService.GetBalance(getBalanceReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeUserWalletNotExist, rsp.Code)
	})

	t.Run("case3: get balance fail-[query db fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		getBalanceReq := &data.GetBalanceReq{UserID: 101}
		mock.ExpectQuery("SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = ? ").WithArgs(getBalanceReq.UserID).WillReturnError(errors.New("db error"))
		walletService := service.NewWalletService(ctx, logID, mockDBCli, nil)
		rsp, err := walletService.GetBalance(getBalanceReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeQueryDBFail, rsp.Code)
	})
}

func TestGetTransactionHistory(t *testing.T) {
	ConnectDBRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectDBRedis()

	mockDBCli := db.GetDbClientMock()
	mock := db.GetSqlMock()
	t.Run("case1: get transaction history success-[history exist]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		page, limit := int32(1), int32(10)
		getHisReq := &data.GetTransactionHistoryReq{UserID: 101, Page: page, Limit: limit}
		tn := time.Now().Unix()
		rows := sqlmock.NewRows([]string{"id", "order_id", "user_id", "tx_type", "amount", "related_user_id", "created_at", "updated_at"})
		rows.AddRow(1, "111", 101, 1, 2000.00, 0, tn, tn)
		rows.AddRow(2, "222", 101, 2, 1000.00, 0, tn, tn)
		rows.AddRow(3, "333", 101, 3, 3000.00, 102, tn, tn)
		rows.AddRow(4, "444", 101, 4, 3000.00, 102, tn, tn)
		offset := (page - 1) * limit
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at FROM transactions WHERE user_id = $1 LIMIT $2 OFFSET $3")).
			WithArgs(getHisReq.UserID, limit, offset).WillReturnRows(rows)
		walletService := service.NewWalletService(ctx, logID, mockDBCli, nil)
		rsp, err := walletService.GetTransactionHistory(getHisReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeSuccess, rsp.Code)
		assert.Equal(t, 4, len(rsp.Data.Items))
		assert.Equal(t, "222", rsp.Data.Items[1].OrderID)
	})

	t.Run("case2: get transaction history success-[history empty]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		page, limit := int32(1), int32(10)
		getHisReq := &data.GetTransactionHistoryReq{UserID: 101, Page: page, Limit: limit}
		offset := (page - 1) * limit
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at FROM transactions WHERE user_id = $1 LIMIT $2 OFFSET $3")).
			WithArgs(getHisReq.UserID, limit, offset).WillReturnError(sql.ErrNoRows)
		walletService := service.NewWalletService(ctx, logID, mockDBCli, nil)
		rsp, err := walletService.GetTransactionHistory(getHisReq)
		assert.Nil(t, err)
		assert.Equal(t, errcode.ErrCodeTransactionNotExist, rsp.Code)
	})

	t.Run("case3: get transaction history fail-[query db fail]", func(t *testing.T) {
		logID := util.Uniqid()
		ctx := context.Background()
		page, limit := int32(1), int32(10)
		getHisReq := &data.GetTransactionHistoryReq{UserID: 101, Page: page, Limit: limit}
		offset := (page - 1) * limit
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at FROM transactions WHERE user_id = $1 LIMIT $2 OFFSET $3")).
			WithArgs(getHisReq.UserID, limit, offset).WillReturnError(errors.New("query db fail"))
		walletService := service.NewWalletService(ctx, logID, mockDBCli, nil)
		rsp, err := walletService.GetTransactionHistory(getHisReq)
		assert.NotNil(t, err)
		assert.Equal(t, errcode.ErrCodeQueryDBFail, rsp.Code)
	})
}
