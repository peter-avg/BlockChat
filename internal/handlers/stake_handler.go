package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Set Stake for Proof of Stake
// ============================
func SetStake(c *gin.Context, MyNode *model.Node) {
	var request model.SetStakeRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if MyNode.Wallet.DeductMoney(request.Stake) == false {
		log.Println("Could not set stake, insufficient funds")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not set stake, insufficient funds"})
		return
	}

	nodeInfo := model.NodeInfo{
		Id:    MyNode.Id,
		Stake: request.Stake,
	}

	if MyNode.BroadcastStake(nodeInfo) {
		log.Println("Stake broadcasted")
		log.Println("Setting stake", request.Stake)
		MyNode.Ring[MyNode.Id].Stake = request.Stake
		c.JSON(http.StatusOK, gin.H{
			"message": "Stake set",
		})
		return
	}

	log.Println("Stake not set")
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Stake not set",
	})
}

// Validate Stake
// ==============
func ValidateStake(c *gin.Context, MyNode *model.Node) {
	var request model.NodeInfo

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Stake > MyNode.Ring[request.Id].Balance {
		log.Println("Insufficient funds to set stake")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient funds to set stake",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stake validated",
	})
}

// Receive Stake
// =============
func ReceiveStake(c *gin.Context, MyNode *model.Node) {
	var request model.NodeInfo

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	MyNode.Ring[request.Id].Stake = request.Stake

	c.JSON(http.StatusOK, gin.H{
		"message": "Stake received",
	})
}
