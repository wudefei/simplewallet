package validator

import (
	"errors"
	"simplewallet/data"
	"simplewallet/util"
)

type ValidatorSvc struct {
}

func NewValidatorSvc() *ValidatorSvc {
	return &ValidatorSvc{}
}
func (v *ValidatorSvc) ValidatorDepositReq(req *data.DepositReq) error {
	if req.OrderID == "" {
		return errors.New("order_id is required")
	}
	if util.CompareFloat(req.Amount, 1e-8, 8) < 0 {
		return errors.New("amount must >= 1e-8")
	}
	if req.UserID <= 0 {
		return errors.New("user_id should > 0")
	}
	return nil
}
func (v *ValidatorSvc) ValidatorWithdrawReq(req *data.WithdrawReq) error {
	if req.OrderID == "" {
		return errors.New("order_id is required")
	}
	if util.CompareFloat(req.Amount, 1e-8, 8) < 0 {
		return errors.New("amount must >= 1e-8")
	}
	if req.UserID <= 0 {
		return errors.New("user_id should > 0")
	}
	return nil
}
func (v *ValidatorSvc) ValidatorTransferReq(req *data.TransferReq) error {
	if req.OrderID == "" {
		return errors.New("order_id is required")
	}
	if util.CompareFloat(req.Amount, 1e-8, 8) < 0 {
		return errors.New("amount must be greater than 0")
	}
	if req.FromUserID <= 0 {
		return errors.New("from_user_id should > 0")
	}
	if req.ToUserID <= 0 {
		return errors.New("to_user_id should > 0")
	}
	if req.FromUserID == req.ToUserID {
		return errors.New("from_user_id and to_user_id must be different")
	}
	return nil
}
func (v *ValidatorSvc) ValidatorGetBalanceReq(req *data.GetBalanceReq) error {
	if req.UserID <= 0 {
		return errors.New("user_id should > 0")
	}
	return nil
}
func (v *ValidatorSvc) ValidatorGetTransactionHistoryReq(req *data.GetTransactionHistoryReq) error {
	if req.UserID <= 0 {
		return errors.New("user_id should > 0")
	}
	if req.Page <= 0 {
		return errors.New("page should > 0")
	}
	if req.Limit <= 0 {
		return errors.New("limit should > 0")
	}
	if req.Limit > 100 {
		return errors.New("limit should <= 100")
	}
	return nil
}
