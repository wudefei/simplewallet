package model

type Transactions struct {
	ID            int64   `db:"id"`
	OrderID       string  `db:"order_id"`
	UserID        int64   `db:"user_id"`
	TxType        int32   `db:"tx_type"`
	Amount        float64 `db:"amount"`
	RelatedUserID int64   `db:"related_user_id"`
	CreatedAt     int64   `db:"created_at"`
	UpdatedAt     int64   `db:"updated_at"`
}
