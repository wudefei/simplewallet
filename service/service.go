package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"simplewallet/data"
	"simplewallet/model"
	"simplewallet/service/dao"
	"simplewallet/util"
	"simplewallet/util/errcode"
	"time"
)

type WalletService struct {
	logID  string
	ctx    context.Context
	dbCli  *sql.DB
	locker util.DistributedLock
}

func NewWalletService(ctx context.Context, logID string, dbCli *sql.DB, locker util.DistributedLock) *WalletService {
	return &WalletService{
		ctx:    ctx,
		logID:  logID,
		dbCli:  dbCli,
		locker: locker,
	}
}

func (s *WalletService) Deposit(req *data.DepositReq) (*data.CommRsp, error) {
	rsp := &data.CommRsp{Code: errcode.ErrCodeInternalErr, Message: errcode.ErrMsgMap[errcode.ErrCodeInternalErr], LogID: s.logID}

	err := s.locker.Lock()
	if err != nil {
		rsp.Code = errcode.ErrCodeLockFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	defer func() {
		if errt := s.locker.UnLock(); errt != nil {
			rsp.Code = errcode.ErrCodeUnLockFail
			rsp.Message = errcode.ErrMsgMap[rsp.Code] + errt.Error()
		}
	}()

	tx, err := s.dbCli.Begin()
	if err != nil {
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}

	// check order_id
	transDao := dao.NewTransactionsDao(s.ctx, s.logID)
	trans, err := transDao.GetTransactionByOrderID(tx, req.OrderID)
	if err != nil {
		_ = tx.Rollback()
		rsp.Code = errcode.ErrCodeQueryDBFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if trans != nil {
		_ = tx.Rollback()
		err = errors.New("order_id already exists")
		rsp.Code = errcode.ErrCodeOrderIDRepeat
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Update balance
	walletDao := dao.NewWalletDao(s.ctx, s.logID)
	err = walletDao.CreateOrUpdateWallet(tx, req.UserID, req.Amount)
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to update balance" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Record transaction
	err = transDao.InsertTransaction(tx, &model.Transactions{OrderID: req.OrderID, UserID: req.UserID, TxType: data.TxTypeDeposit, Amount: req.Amount})
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to record transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	rsp.Code = errcode.ErrCodeSuccess
	rsp.Message = "Deposit successful"
	return rsp, nil
}

func (s *WalletService) Withdraw(req *data.WithdrawReq) (*data.CommRsp, error) {
	rsp := &data.CommRsp{Code: errcode.ErrCodeInternalErr, Message: errcode.ErrMsgMap[errcode.ErrCodeInternalErr], LogID: s.logID}

	err := s.locker.Lock()
	if err != nil {
		rsp.Code = errcode.ErrCodeLockFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	defer func() {
		if errt := s.locker.UnLock(); errt != nil {
			rsp.Code = errcode.ErrCodeUnLockFail
			rsp.Message = errcode.ErrMsgMap[rsp.Code] + errt.Error()
		}
	}()

	tx, err := s.dbCli.Begin()
	if err != nil {
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	// check order_id
	transDao := dao.NewTransactionsDao(s.ctx, s.logID)
	trans, err := transDao.GetTransactionByOrderID(tx, req.OrderID)
	if err != nil {
		_ = tx.Rollback()
		rsp.Code = errcode.ErrCodeQueryDBFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if trans != nil {
		_ = tx.Rollback()
		err = errors.New("order_id already exists")
		rsp.Code = errcode.ErrCodeOrderIDRepeat
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Check balance
	walletDao := dao.NewWalletDao(s.ctx, s.logID)
	wallet, err := walletDao.GetWalletByUserID(nil, tx, req.UserID)
	if err != nil {
		_ = tx.Rollback()
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if wallet == nil {
		_ = tx.Rollback()
		err = errors.New("user wallet not exist")
		rsp.Code = errcode.ErrCodeUserWalletNotExist
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Compare balance[float64] with withdraw amount[float64]
	if util.CompareFloat(wallet.Balance, req.Amount, 8) < 0 {
		_ = tx.Rollback()
		err = errors.New("balance not enough")
		rsp.Code = errcode.ErrCodeBalanceNotEnough
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Update balance
	err = walletDao.UpdateWalletBalance(tx, req.UserID, data.TxTypeWithdraw, req.Amount)
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to update balance" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Record transaction
	err = transDao.InsertTransaction(tx, &model.Transactions{OrderID: req.OrderID, UserID: req.UserID, TxType: data.TxTypeWithdraw, Amount: req.Amount})
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to record transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	rsp.Code = errcode.ErrCodeSuccess
	rsp.Message = "Withdrawal successful"
	return rsp, nil
}

func (s *WalletService) Transfer(req *data.TransferReq) (*data.CommRsp, error) {
	rsp := &data.CommRsp{Code: errcode.ErrCodeInternalErr, Message: errcode.ErrMsgMap[errcode.ErrCodeInternalErr], LogID: s.logID}

	err := s.locker.Lock()
	if err != nil {
		rsp.Code = errcode.ErrCodeLockFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	defer func() {
		if errt := s.locker.UnLock(); errt != nil {
			rsp.Code = errcode.ErrCodeUnLockFail
			rsp.Message = errcode.ErrMsgMap[rsp.Code] + errt.Error()
		}
	}()

	tx, err := s.dbCli.Begin()
	if err != nil {
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	// check order_id
	transDao := dao.NewTransactionsDao(s.ctx, s.logID)
	trans, err := transDao.GetTransactionByOrderID(tx, req.OrderID)
	if err != nil {
		_ = tx.Rollback()
		rsp.Code = errcode.ErrCodeQueryDBFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if trans != nil {
		_ = tx.Rollback()
		err = errors.New("order_id already exists")
		rsp.Code = errcode.ErrCodeOrderIDRepeat
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Check balance of sender
	walletDao := dao.NewWalletDao(s.ctx, s.logID)
	wallet, err := walletDao.GetWalletByUserID(nil, tx, req.FromUserID)
	if err != nil {
		_ = tx.Rollback()
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if wallet == nil {
		_ = tx.Rollback()
		err = errors.New("sender wallet not exist")
		rsp.Code = errcode.ErrCodeUserWalletNotExist
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Compare balance[float64] with transfer amount[float64]
	if util.CompareFloat(wallet.Balance, req.Amount, 8) < 0 {
		_ = tx.Rollback()
		err = errors.New("balance not enough")
		rsp.Code = errcode.ErrCodeBalanceNotEnough
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Update sender's balance
	err = walletDao.UpdateWalletBalance(tx, req.FromUserID, data.TxTypeTransferOut, req.Amount)
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to update sender's balance" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Update recipient's balance
	err = walletDao.CreateOrUpdateWallet(tx, req.ToUserID, req.Amount)
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to update recipient's balance" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	// Record transactions
	err = transDao.InsertTransaction(tx, &model.Transactions{OrderID: req.OrderID, UserID: req.FromUserID, TxType: data.TxTypeTransferOut, Amount: req.Amount, RelatedUserID: req.ToUserID})
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to record sender's transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}
	err = transDao.InsertTransaction(tx, &model.Transactions{OrderID: req.OrderID, UserID: req.ToUserID, TxType: data.TxTypeTransferIn, Amount: req.Amount, RelatedUserID: req.FromUserID})
	if err != nil {
		_ = tx.Rollback()
		log.Println("Failed to record recipient's transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit transaction" + err.Error())
		rsp.Code = errcode.ErrCodeDbError
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}

	rsp.Code = errcode.ErrCodeSuccess
	rsp.Message = "Transfer successful"
	return rsp, nil
}

func (s *WalletService) GetBalance(req *data.GetBalanceReq) (*data.GetBalanceRsp, error) {
	rsp := &data.GetBalanceRsp{Code: errcode.ErrCodeInternalErr, Message: errcode.ErrMsgMap[errcode.ErrCodeInternalErr], LogID: s.logID}

	wallet, err := dao.NewWalletDao(s.ctx, s.logID).GetWalletByUserID(s.dbCli, nil, req.UserID)
	if err != nil {
		log.Println("Failed to get balance" + err.Error())
		rsp.Code = errcode.ErrCodeQueryDBFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if wallet == nil {
		err = errors.New("wallet not exist")
		rsp.Code = errcode.ErrCodeUserWalletNotExist
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, err
	}
	rsp.Code = errcode.ErrCodeSuccess
	rsp.Message = "Success"
	rsp.Data = &data.GetBalanceRspData{Balance: wallet.Balance}
	return rsp, nil
}

func (s *WalletService) GetTransactionHistory(req *data.GetTransactionHistoryReq) (*data.GetTransactionHistoryRsp, error) {
	rsp := &data.GetTransactionHistoryRsp{Code: errcode.ErrCodeInternalErr, Message: errcode.ErrMsgMap[errcode.ErrCodeInternalErr], LogID: s.logID}
	var rspItems []*data.GetTransactionHistoryRspDataItem

	transDao := dao.NewTransactionsDao(s.ctx, s.logID)
	txList, err := transDao.GetTransactionListByUserID(s.dbCli, req.UserID, req.Page, req.Limit)
	if err != nil {
		log.Println("Failed to get transaction history" + err.Error())
		rsp.Code = errcode.ErrCodeQueryDBFail
		rsp.Message = errcode.ErrMsgMap[rsp.Code] + err.Error()
		return rsp, err
	}
	if txList == nil {
		rsp.Code = errcode.ErrCodeTransactionNotExist
		rsp.Message = errcode.ErrMsgMap[rsp.Code]
		return rsp, nil
	}
	for _, tx := range txList {
		rspItems = append(rspItems, &data.GetTransactionHistoryRspDataItem{
			OrderID:       tx.OrderID,
			UserID:        tx.UserID,
			TxType:        tx.TxType,
			Amount:        tx.Amount,
			RelatedUserID: tx.RelatedUserID,
			CreatedAt:     time.Unix(tx.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		})
	}

	rsp.Code = errcode.ErrCodeSuccess
	rsp.Message = "Success"
	rsp.Data = &data.GetTransactionHistoryRspData{Items: rspItems}
	return rsp, nil
}
