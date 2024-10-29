package dao

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"simplewallet/model"
	"time"
)

type TransactionsDao struct {
	ctx   context.Context
	logID string
}

func NewTransactionsDao(ctx context.Context, logID string) *TransactionsDao {
	return &TransactionsDao{ctx: ctx, logID: logID}
}

func (d *TransactionsDao) GetTransactionByOrderID(dbTx *sql.Tx, orderID string) (*model.Transactions, error) {
	tx := &model.Transactions{}
	err := dbTx.QueryRow("SELECT id,order_id,user_id,tx_type,amount,related_user_id,created_at,updated_at FROM transactions WHERE order_id = $1", orderID).
		Scan(&tx.ID, &tx.OrderID, &tx.UserID, &tx.TxType, &tx.Amount, &tx.RelatedUserID, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("%s|[%s] Failed to get transactions by order id: %v", d.logID, orderID, err)
		return nil, err
	}
	return tx, nil
}

func (d *TransactionsDao) GetTransactionListByUserID(db *sql.DB, userID int64, page int32, limit int32) ([]*model.Transactions, error) {
	txList := make([]*model.Transactions, 0)
	offset := (page - 1) * limit
	rows, err := db.Query("SELECT id, order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at FROM transactions WHERE user_id = $1 LIMIT $2 OFFSET $3", userID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("%s|[%d] Failed to get transaction list by user_id: %v", d.logID, userID, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tx := &model.Transactions{}
		if err = rows.Scan(&tx.ID, &tx.OrderID, &tx.UserID, &tx.TxType, &tx.Amount, &tx.RelatedUserID, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			log.Printf("%s|[%d] Failed to scan transaction: %v", d.logID, userID, err)
			return nil, err
		}
		txList = append(txList, tx)
	}
	return txList, nil
}

func (d *TransactionsDao) InsertTransaction(dbTx *sql.Tx, tx *model.Transactions) error {
	tn := time.Now().Unix()
	_, err := dbTx.Exec("INSERT INTO transactions (order_id, user_id, tx_type, amount, related_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		tx.OrderID, tx.UserID, tx.TxType, tx.Amount, tx.RelatedUserID, tn, tn)
	return err
}
