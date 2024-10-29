package data

type DepositReq struct {
	OrderID string  `json:"order_id"`
	UserID  int64   `json:"user_id"`
	Amount  float64 `json:"amount"`
}
type CommRsp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	LogID   string `json:"log_id"`
}

type WithdrawReq struct {
	OrderID string  `json:"order_id"`
	UserID  int64   `json:"user_id"`
	Amount  float64 `json:"amount"`
}

type TransferReq struct {
	OrderID    string  `json:"order_id"`
	FromUserID int64   `json:"from_user_id"`
	ToUserID   int64   `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

type GetBalanceReq struct {
	UserID int64 `json:"user_id"`
}
type GetBalanceRsp struct {
	Code    int32              `json:"code"`
	Message string             `json:"message"`
	Data    *GetBalanceRspData `json:"data"`
	LogID   string             `json:"log_id"`
}
type GetBalanceRspData struct {
	Balance float64 `json:"balance"`
}

type GetTransactionHistoryReq struct {
	UserID int64 `json:"user_id"`
	Page   int32 `json:"page"`
	Limit  int32 `json:"limit"`
}
type GetTransactionHistoryRsp struct {
	Code    int32                         `json:"code"`
	Message string                        `json:"message"`
	Data    *GetTransactionHistoryRspData `json:"data"`
	LogID   string                        `json:"log_id"`
}
type GetTransactionHistoryRspData struct {
	Items []*GetTransactionHistoryRspDataItem `json:"items"`
}
type GetTransactionHistoryRspDataItem struct {
	OrderID       string  `json:"order_id"`
	UserID        int64   `json:"user_id"`
	TxType        int32   `json:"tx_type"` //0: undefine 1: deposit, 2: withdrawal, 3: transfer in, 4: transfer out
	Amount        float64 `json:"amount"`
	RelatedUserID int64   `json:"related_user_id"`
	CreatedAt     string  `json:"created_at"`
}
