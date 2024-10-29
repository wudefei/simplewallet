package router

import (
	"net/http"
	"simplewallet/controller"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    http.StatusNotFound,
			"message": "Not exists",
			"data":    map[string]string{},
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    http.StatusMethodNotAllowed,
			"message": "Fobidden method",
			"data":    map[string]string{},
		})
	})

	router.Any("/", func(c *gin.Context) {
		c.Header("server", "http/1.1")
		c.JSON(200, gin.H{
			"message": "service " + time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	ctl := controller.NewWalletController()
	api := router.Group("")
	{
		api.POST("/deposit", ctl.Deposit)
		api.POST("/withdraw", ctl.Withdraw)
		api.POST("/transfer", ctl.Transfer)
		api.GET("/balance", ctl.GetBalance)
		api.GET("/transactions", ctl.GetTransactionHistory)
	}

	return router
}
