package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Get Balance from Wallet
// =======================
func GetBalance(c *gin.Context, MyNode *model.Node) {
	var balance float64
	var softStake float64
	for _, nodeInfo := range MyNode.Ring {
		if nodeInfo.Id == MyNode.Id {
			balance = nodeInfo.SoftBalance
			softStake = nodeInfo.SoftStake
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"balance":    balance,
		"soft_stake": softStake,
	})
}
