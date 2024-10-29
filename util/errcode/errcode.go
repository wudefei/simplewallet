package errcode

const (
	ErrCodeSuccess             int32 = 0
	ErrCodeFail                int32 = 1000
	ErrCodeBadRequestParam     int32 = 1001
	ErrCodeDbError             int32 = 1002
	ErrCodeLockFail            int32 = 1003
	ErrCodeUnLockFail          int32 = 1004
	ErrCodeUserWalletNotExist  int32 = 1005
	ErrCodeOrderIDRepeat       int32 = 1006
	ErrCodeBalanceNotEnough    int32 = 1007
	ErrCodeQueryDBFail         int32 = 1008
	ErrCodeInternalErr         int32 = 1009
	ErrCodeTransactionNotExist int32 = 1010
)

var (
	ErrMsgMap = map[int32]string{
		ErrCodeSuccess:             "success",
		ErrCodeFail:                "fail",
		ErrCodeBadRequestParam:     "bad request param",
		ErrCodeDbError:             "db error",
		ErrCodeLockFail:            "lock fail",
		ErrCodeUnLockFail:          "unlock fail",
		ErrCodeUserWalletNotExist:  "user wallet not exist",
		ErrCodeOrderIDRepeat:       "order_id repeat",
		ErrCodeBalanceNotEnough:    "insufficient funds",
		ErrCodeQueryDBFail:         "failed to query db",
		ErrCodeInternalErr:         "internal error",
		ErrCodeTransactionNotExist: "transaction not exist",
	}
)
