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
func SetStake(c *gin.Context, myNode *model.Node) {
	// Stake is a Transaction with a Recipient Address == -1
	var request model.SetStakeRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stakeAmount = request.Stake
	var receiverAddress = &config.STAKE_PUBLIC_ADDRESS

	var stakeTransaction = model.NewTransaction(
		receiverAddress,
		true,
		strconv.FormatFloat(stakeAmount, 'f', -1, 64),
		myNode.Wallet.AddTransaction(),
	)

	var err error
	stakeTransaction.Signature, err = myNode.Wallet.SignTransaction(stakeTransaction)
	stakeTransaction.SenderAddress = myNode.Wallet.PublicKey
	if err != nil {
		log.Println("Error signing stake transaction", err)
	}

	for _, node := range myNode.Ring {
		if node.Id == myNode.Id {
			senderBalance := node.SoftBalance - node.SoftStake
			if senderBalance < stakeAmount {
				log.Println("Could not set stake, insufficient funds")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Could not set stake, insufficient funds",
				})
				return
			}
		}
	}
	//isBlockFull := myNode.CurrentBlock.AddTransaction(*stakeTransaction, myNode)
	txnPoolObj := model.TransactionInPool{Txn: *stakeTransaction, IsSendTxn: true}
	model.TransactionPool <- txnPoolObj
	if myNode.BroadcastTransaction(stakeTransaction) {
		log.Println("Stake Transaction broadcasted")
		c.JSON(http.StatusOK, gin.H{
			"message": "Stake Transaction of amount " + strconv.FormatFloat(stakeAmount, 'f', -1, 64) + " broadcasted",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stake transaction not sent",
		})
	}

	//if isBlockFull {
	//	myNode.CurrentBlock.ElectLeader(myNode)
	//}

}
