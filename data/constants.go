package data

const LogIdParam string = "logId"

// 0: undefine 1: deposit, 2: withdrawal, 3: transfer in, 4: transfer out
const (
	TxTypeUnknown     int32 = 0
	TxTypeDeposit     int32 = 1
	TxTypeWithdraw    int32 = 2
	TxTypeTransferIn  int32 = 3
	TxTypeTransferOut int32 = 4
)
