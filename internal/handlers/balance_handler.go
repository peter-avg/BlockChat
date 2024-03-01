package handlers

import (
	"block-chat/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Get Balance from Wallet
// =======================
func GetBalance(c *gin.Context, MyNode *model.Node) {
	c.JSON(http.StatusOK, gin.H{
		"balance": fmt.Sprint(MyNode.Wallet.Balance),
	})
}
