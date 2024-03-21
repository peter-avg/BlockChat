package handlers

import (
	"block-chat/internal/config"
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// Broadcast a stake transaction
// ===================================
func SetStake(c *gin.Context, MyNode *model.Node) {
	// Stake is a Transaction with a Recipient Address == -1

	var request model.SetStakeRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stakeAmount int = request.Stake
	var receiverAddress int = -1

	var stakeTransaction = model.NewTransaction(
		receiverAddress,
		true,
		"",
		MyNode.Wallet.AddTransaction(),
	)

	var err error
	stakeTransaction.Signature, err = MyNode.Wallet.SignTransaction(stakeTransaction)
	stakeTransaction.SenderAddress = MyNode.Wallet.PublicKey
	if err != nil {
		log.Println("Error signing stake transaction", err)
	}

	if MyNode.Wallet.DeductMoney(stakeAmount) == false {
		log.Println("Could not set stake, insufficient funds")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not set stake, insufficient funds"})
		return
	}

	if MyNode.BroadcastTransaction(stakeTransaction) {
		log.Println("Stake Transaction broadcasted")
		MyNode.CurrentBlock.AddTransaction(*stakeTransaction, config.CAPACITY)
		c.JSON(http.StatusOK, gin.H{
			"message": "Stake Transaction of amount " + strconv.Itoa(stakeAmount) + " broadcasted",
		})

		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Stake transaction not sent",
	})

	//
	//nodeInfo := model.NodeInfo{
	//	Id:    MyNode.Id,
	//	Stake: request.Stake,
	//}
	//
	//if MyNode.BroadcastStake(nodeInfo) {
	//	log.Println("Stake broadcasted")
	//	log.Println("Setting stake", request.Stake)
	//	MyNode.Ring[MyNode.Id].Stake = request.Stake
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "Stake set",
	//	})
	//	return
	//}
	//
	//log.Println("Stake not set")
	//c.JSON(http.StatusBadRequest, gin.H{
	//	"error": "Stake not set",
	//})
}
