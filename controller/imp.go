package controller

import (
	"errors"
	"log"
	"net/http"
	"simplewallet/controller/validator"
	"simplewallet/data"
	"simplewallet/service"
	"simplewallet/util"
	"simplewallet/util/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
}

func NewWalletController() *WalletController {
	return &WalletController{}
}
func (w *WalletController) Deposit(ctx *gin.Context) {
	logID := util.Uniqid()
	defer func() {
		if p := recover(); p != nil {
			log.Printf("%s|panic:%v\n", logID, p)
		}
	}()
	var req data.DepositReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.NewValidatorSvc().ValidatorDepositReq(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbCli := db.GetDbClient()
	lockKey := "deposit:" + req.OrderID
	locker := util.NewDistributedLock(logID, lockKey, 5)

	s := service.NewWalletService(ctx, logID, dbCli, locker)
	rsp, err := s.Deposit(&req)
	if err != nil {
		log.Printf("%s|fail to deposit:%s\n", logID, err.Error())
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (w *WalletController) Withdraw(ctx *gin.Context) {
	logID := util.Uniqid()
	defer func() {
		if p := recover(); p != nil {
			log.Printf("%s|panic:%v\n", logID, p)
		}
	}()
	var req data.WithdrawReq
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = validator.NewValidatorSvc().ValidatorWithdrawReq(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbCli := db.GetDbClient()
	lockKey := "withdraw:" + req.OrderID
	locker := util.NewDistributedLock(logID, lockKey, 5)

	s := service.NewWalletService(ctx, logID, dbCli, locker)
	rsp, err := s.Withdraw(&req)
	if err != nil {
		log.Printf("%s|fail to withdraw:%s\n", logID, err.Error())
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (w *WalletController) Transfer(ctx *gin.Context) {
	logID := util.Uniqid()
	defer func() {
		if p := recover(); p != nil {
			log.Printf("%s|panic:%v\n", logID, p)
		}
	}()
	var req data.TransferReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validator.NewValidatorSvc().ValidatorTransferReq(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbCli := db.GetDbClient()
	lockKey := "transfer:" + req.OrderID
	locker := util.NewDistributedLock(logID, lockKey, 5)

	s := service.NewWalletService(ctx, logID, dbCli, locker)
	rsp, err := s.Transfer(&req)
	if err != nil {
		log.Printf("%s|fail to transfer:%s\n", logID, err.Error())
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (w *WalletController) GetBalance(ctx *gin.Context) {
	logID := util.Uniqid()
	defer func() {
		if p := recover(); p != nil {
			log.Printf("%s|panic:%v\n", logID, p)
		}
	}()
	userID, err := w.GetParamUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := &data.GetBalanceReq{UserID: userID}
	if err := validator.NewValidatorSvc().ValidatorGetBalanceReq(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbCli := db.GetDbClient()
	s := service.NewWalletService(ctx, logID, dbCli, nil)
	rsp, err := s.GetBalance(req)
	if err != nil {
		log.Printf("%s|fail to get balance:%s\n", logID, err.Error())
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (w *WalletController) GetTransactionHistory(ctx *gin.Context) {
	logID := util.Uniqid()
	defer func() {
		if p := recover(); p != nil {
			log.Printf("%s|panic:%v\n", logID, p)
		}
	}()
	userID, err := w.GetParamUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := w.GetParamPage(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, err := w.GetParamLimit(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &data.GetTransactionHistoryReq{UserID: userID, Page: page, Limit: limit}
	if err := validator.NewValidatorSvc().ValidatorGetTransactionHistoryReq(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbCli := db.GetDbClient()
	s := service.NewWalletService(ctx, logID, dbCli, nil)
	rsp, err := s.GetTransactionHistory(req)
	if err != nil {
		log.Printf("%s|fail to get transaction history:%s\n", logID, err.Error())
	}
	ctx.JSON(http.StatusOK, rsp)
}
func (w *WalletController) GetParamUserID(ctx *gin.Context) (int64, error) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		return 0, errors.New("user_id is required")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
func (w *WalletController) GetParamPage(ctx *gin.Context) (int32, error) {
	pageStr := ctx.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}
	return int32(page), nil
}
func (w *WalletController) GetParamLimit(ctx *gin.Context) (int32, error) {
	limitStr := ctx.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}
	return int32(limit), nil
}
