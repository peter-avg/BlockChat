package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Get Balance from Wallet
// =======================
func GetBalance(c *gin.Context, MyNode *model.Node) {
	c.JSON(http.StatusOK, gin.H{
		"balance": MyNode.Wallet.Balance,
	})
}
