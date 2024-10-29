package model

type Wallet struct {
	ID        int64   `db:"id"`
	UserID    int64   `db:"user_id"`
	Balance   float64 `db:"balance"`
	CreatedAt int64   `db:"created_at"`
	UpdatedAt int64   `db:"updated_at"`
}
