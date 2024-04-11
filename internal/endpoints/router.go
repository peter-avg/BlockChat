package endpoints

import (
	"block-chat/internal/handlers"
	"block-chat/internal/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(myNode *model.Node) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// Bootstrap Endpoints of API
	// ==========================
	router.POST("/blockchat_api/register_node", func(c *gin.Context) {
		handlers.RegisterNode(c, myNode)
	})

	// Backend Endpoints of API
	// ========================
	router.POST("/blockchat_api/receive_new_node", func(c *gin.Context) {
		handlers.ReceiveNewNode(c, myNode)
	})
	router.POST("/blockchat_api/validate_transaction", func(c *gin.Context) {
		handlers.ValidateTransaction(c, myNode)
	})
	router.POST("/blockchat_api/receive_validated_block", func(c *gin.Context) {
		handlers.ReceiveValidatedBlock(c, myNode)
	})

	// Client Endpoints of API
	// =======================
	router.POST("/blockchat_api/set_stake", func(c *gin.Context) {
		handlers.SetStake(c, myNode)
	})
	router.POST("/blockchat_api/send_transaction", func(c *gin.Context) {
		handlers.SendTransaction(c, myNode)
	})
	router.GET("/blockchat_api/get_balance", func(c *gin.Context) {
		handlers.GetBalance(c, myNode)
	})
	router.GET("/blockchat_api/get_last_block", func(c *gin.Context) {
		handlers.GetLastBlock(c, myNode)
	})

	go processTransactions(myNode)

	return router
}
