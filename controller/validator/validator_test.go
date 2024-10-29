package validator_test

import (
	"errors"
	"simplewallet/controller/validator"
	"simplewallet/data"
	"testing"

	"go.uber.org/goleak"
)

func TestValidatorDepositReq(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks

	type args struct {
		Name string
		args *data.DepositReq
		want error
	}
	tests := []args{
		{Name: "case1: ValidatorDepositReq success", args: &data.DepositReq{OrderID: "123", UserID: 101, Amount: 1000.00}, want: nil},
		{Name: "case2: ValidatorDepositReq fail-[order_id is empty]", args: &data.DepositReq{OrderID: "", UserID: 101, Amount: 1000.00}, want: errors.New("order_id is required")},
		{Name: "case3: ValidatorDepositReq fail-[UserID = 0]", args: &data.DepositReq{OrderID: "123", UserID: 0, Amount: 1000.00}, want: errors.New("user_id is required")},
		{Name: "case4: ValidatorDepositReq fail-[UserID < 0]", args: &data.DepositReq{OrderID: "123", UserID: -101, Amount: 1000.00}, want: errors.New("user_id should > 0")},
		{Name: "case5: ValidatorDepositReq fail-[amount = 0]", args: &data.DepositReq{OrderID: "123", UserID: 101, Amount: 0}, want: errors.New("amount should > 0")},
		{Name: "case6: ValidatorDepositReq fail-[amount < 0]", args: &data.DepositReq{OrderID: "123", UserID: 101, Amount: -1000.00}, want: errors.New("amount should > 0")},
		{Name: "case7: ValidatorDepositReq success-[amount = 0.00000001]", args: &data.DepositReq{OrderID: "123", UserID: 101, Amount: 0.00000001}, want: nil},
		{Name: "case8: ValidatorDepositReq fail-   [amount =-0.00000001]", args: &data.DepositReq{OrderID: "123", UserID: 101, Amount: -0.00000001}, want: errors.New("amount should > 0")},
		{Name: "case9: ValidatorDepositReq fail-   [amount = 0.000000001]", args: &data.DepositReq{OrderID: "123", UserID: 101, Amount: 0.000000001}, want: errors.New("amount should >= 1e-8")},
	}
	v := validator.NewValidatorSvc()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.ValidatorDepositReq(tt.args)
			if tt.want == nil && err != nil {
				t.Errorf("ValidatorDepositReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			} else if tt.want != nil && err == nil {
				t.Errorf("ValidatorDepositReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			}
		})
	}
}

func TestValidatorWithdrawReq(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks
	type args struct {
		Name string
		args *data.WithdrawReq
		want error
	}
	tests := []args{
		{Name: "case1: ValidatorWithdrawReq success", args: &data.WithdrawReq{OrderID: "123", UserID: 101, Amount: 1000.00}, want: nil},
		{Name: "case2: ValidatorWithdrawReq fail-[order_id is empty]", args: &data.WithdrawReq{OrderID: "", UserID: 101, Amount: 1000.00}, want: errors.New("order_id is required")},
		{Name: "case3: ValidatorWithdrawReq fail-[UserID = 0]", args: &data.WithdrawReq{OrderID: "123", UserID: 0, Amount: 1000.00}, want: errors.New("user_id is required")},
		{Name: "case4: ValidatorWithdrawReq fail-[UserID < 0]", args: &data.WithdrawReq{OrderID: "123", UserID: -101, Amount: 1000.00}, want: errors.New("user_id should > 0")},
		{Name: "case5: ValidatorWithdrawReq fail-[amount = 0]", args: &data.WithdrawReq{OrderID: "123", UserID: 101, Amount: 0}, want: errors.New("amount should > 0")},
		{Name: "case6: ValidatorWithdrawReq fail-[amount < 0]", args: &data.WithdrawReq{OrderID: "123", UserID: 101, Amount: -1000.00}, want: errors.New("amount should > 0")},
		{Name: "case7: ValidatorWithdrawReq success-[amount = 0.00000001]", args: &data.WithdrawReq{OrderID: "123", UserID: 101, Amount: 0.00000001}, want: nil},
		{Name: "case8: ValidatorWithdrawReq fail-   [amount =-0.00000001]", args: &data.WithdrawReq{OrderID: "123", UserID: 101, Amount: -0.00000001}, want: errors.New("amount should > 0")},
		{Name: "case9: ValidatorWithdrawReq fail-   [amount = 0.000000001]", args: &data.WithdrawReq{OrderID: "123", UserID: 101, Amount: 0.000000001}, want: errors.New("amount should >= 1e-8")},
	}
	v := validator.NewValidatorSvc()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.ValidatorWithdrawReq(tt.args)
			if tt.want == nil && err != nil {
				t.Errorf("ValidatorWithdrawReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			} else if tt.want != nil && err == nil {
				t.Errorf("ValidatorWithdrawReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			}
		})
	}
}

func TestValidatorTransferReq(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks
	type args struct {
		Name string
		args *data.TransferReq
		want error
	}
	tests := []args{
		{Name: "case1: ValidatorTransferReq success", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 102, Amount: 1000.00}, want: nil},
		{Name: "case2: ValidatorTransferReq fail-[order_id is empty]", args: &data.TransferReq{OrderID: "", FromUserID: 101, ToUserID: 102, Amount: 1000.00}, want: errors.New("order_id is required")},
		{Name: "case3: ValidatorTransferReq fail-[FromUserID = 0]", args: &data.TransferReq{OrderID: "123", FromUserID: 0, ToUserID: 102, Amount: 1000.00}, want: errors.New("from_user_id should > 0")},
		{Name: "case4: ValidatorTransferReq fail-[ToUserID = 0]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 0, Amount: 1000.00}, want: errors.New("to_user_id should > 0")},
		{Name: "case5: ValidatorTransferReq fail-[FromUserID < 0]", args: &data.TransferReq{OrderID: "123", FromUserID: -101, ToUserID: 102, Amount: 1000.00}, want: errors.New("from_user_id should > 0")},
		{Name: "case6: ValidatorTransferReq fail-[ToUserID < 0]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: -102, Amount: 1000.00}, want: errors.New("to_user_id should > 0")},
		{Name: "case7: ValidatorTransferReq fail-[FromUserID = ToUserID]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 101, Amount: 1000.00}, want: errors.New("from_user_id and to_user_id must be different")},
		{Name: "case8: ValidatorTransferReq fail-[amount = 0]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 102, Amount: 0}, want: errors.New("amount should > 0")},
		{Name: "case9: ValidatorTransferReq fail-[amount < 0]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 102, Amount: -1000.00}, want: errors.New("amount should > 0")},
		{Name: "case10: ValidatorTransferReq success-[amount = 0.00000001]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 102, Amount: 0.00000001}, want: nil},
		{Name: "case11: ValidatorTransferReq fail-   [amount =-0.00000001]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 102, Amount: -0.00000001}, want: errors.New("amount should > 0")},
		{Name: "case12: ValidatorTransferReq fail-   [amount = 0.000000001]", args: &data.TransferReq{OrderID: "123", FromUserID: 101, ToUserID: 102, Amount: 0.000000001}, want: errors.New("amount should >= 1e-8")},
	}
	v := validator.NewValidatorSvc()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.ValidatorTransferReq(tt.args)
			if tt.want == nil && err != nil {
				t.Errorf("ValidatorWithdrawReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			} else if tt.want != nil && err == nil {
				t.Errorf("ValidatorWithdrawReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			}
		})
	}
}

func TestValidatorGetBalanceReq(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks
	type args struct {
		Name string
		args *data.GetBalanceReq
		want error
	}
	tests := []args{
		{Name: "case1: ValidatorGetBalanceReq success", args: &data.GetBalanceReq{UserID: 101}, want: nil},
		{Name: "case2: ValidatorGetBalanceReq fail-[UserID = 0]", args: &data.GetBalanceReq{UserID: 0}, want: errors.New("user_id should > 0")},
		{Name: "case3: ValidatorGetBalanceReq fail-[UserID < 0]", args: &data.GetBalanceReq{UserID: -101}, want: errors.New("user_id should > 0")},
	}
	v := validator.NewValidatorSvc()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.ValidatorGetBalanceReq(tt.args)
			if tt.want == nil && err != nil {
				t.Errorf("ValidatorGetBalanceReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			} else if tt.want != nil && err == nil {
				t.Errorf("ValidatorGetBalanceReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			}
		})
	}
}

func TestValidatorGetTransactionHistoryReq(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks

	type args struct {
		Name string
		args *data.GetTransactionHistoryReq
		want error
	}
	tests := []args{
		{Name: "case1: GetTransactionHistoryReq success", args: &data.GetTransactionHistoryReq{UserID: 101, Page: 1, Limit: 10}, want: nil},
		{Name: "case2: GetTransactionHistoryReq fail-[UserID = 0]", args: &data.GetTransactionHistoryReq{UserID: 0, Page: 1, Limit: 10}, want: errors.New("user_id should > 0")},
		{Name: "case3: GetTransactionHistoryReq fail-[UserID < 0]", args: &data.GetTransactionHistoryReq{UserID: -101, Page: 1, Limit: 10}, want: errors.New("user_id should > 0")},
		{Name: "case3: GetTransactionHistoryReq fail-[page = 0]", args: &data.GetTransactionHistoryReq{UserID: 101, Page: 0, Limit: 10}, want: errors.New("page should > 0")},
		{Name: "case3: GetTransactionHistoryReq fail-[page < 0]", args: &data.GetTransactionHistoryReq{UserID: 101, Page: -1, Limit: 10}, want: errors.New("page should > 0")},
		{Name: "case3: GetTransactionHistoryReq fail-[limit = 0]", args: &data.GetTransactionHistoryReq{UserID: 101, Page: 0, Limit: 10}, want: errors.New("limit should > 0")},
		{Name: "case3: GetTransactionHistoryReq fail-[limit < 0]", args: &data.GetTransactionHistoryReq{UserID: 101, Page: -1, Limit: 10}, want: errors.New("limit should > 0")},
		{Name: "case3: GetTransactionHistoryReq fail-[limit > 100]", args: &data.GetTransactionHistoryReq{UserID: 101, Page: -1, Limit: 10}, want: errors.New("limit should <= 100")},
	}
	v := validator.NewValidatorSvc()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := v.ValidatorGetTransactionHistoryReq(tt.args)
			if tt.want == nil && err != nil {
				t.Errorf("ValidatorGetTransactionHistoryReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			} else if tt.want != nil && err == nil {
				t.Errorf("ValidatorGetTransactionHistoryReq()case:%s error = %v, wantErr %v", tt.Name, err, tt.want)
			}
		})
	}
}
