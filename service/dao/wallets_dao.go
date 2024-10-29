package dao

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"simplewallet/data"
	"simplewallet/model"
	"time"
)

type WalletDao struct {
	ctx   context.Context
	logID string
}

func NewWalletDao(ctx context.Context, logID string) *WalletDao {
	return &WalletDao{ctx: ctx, logID: logID}
}

func (d *WalletDao) GetWalletByUserID(db *sql.DB, dbTx *sql.Tx, userID int64) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	querySql := "SELECT id,user_id,balance,created_at,updated_at FROM wallets WHERE user_id = $1"
	var sqlRow *sql.Row
	if dbTx != nil {
		sqlRow = dbTx.QueryRow(querySql, userID)
	} else {
		sqlRow = db.QueryRow(querySql, userID)
	}
	err := sqlRow.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("%s|[%d] Failed to get wallet by user id: %v", d.logID, userID, err)
		return nil, err
	}
	return wallet, nil
}

// create or update
func (d *WalletDao) CreateOrUpdateWallet(dbTx *sql.Tx, userID int64, balance float64) error {
	tn := time.Now().Unix()
	wallet, err := d.GetWalletByUserID(nil, dbTx, userID)
	if err != nil {
		log.Printf("%s|[%d] Failed to query wallet: %v", d.logID, userID, err)
		return err
	}
	if wallet == nil {
		_, err = dbTx.Exec("INSERT INTO wallets (user_id, balance,created_at,updated_at) VALUES ($1, $2, $3, $4)", userID, balance, tn, tn)
	} else {
		_, err = dbTx.Exec("UPDATE wallets SET balance = wallets.balance + $1, updated_at = $2 WHERE user_id = $3", balance, tn, userID)
	}
	if err != nil {
		log.Printf("%s|[%d] Failed to create or update wallet: %v", d.logID, userID, err)
		return err
	}
	return nil
}

// update wallet balance
func (d *WalletDao) UpdateWalletBalance(dbTx *sql.Tx, userID int64, txType int32, balance float64) error {
	tn := time.Now().Unix()
	var err error
	if txType == data.TxTypeDeposit || txType == data.TxTypeTransferIn {
		_, err = dbTx.Exec("UPDATE wallets SET balance = wallets.balance + $1, updated_at = $2 WHERE user_id = $3", balance, tn, userID)
	} else if txType == data.TxTypeWithdraw || txType == data.TxTypeTransferOut {
		_, err = dbTx.Exec("UPDATE wallets SET balance = wallets.balance - $1, updated_at = $2 WHERE user_id = $3", balance, tn, userID)
	}
	return err
}
