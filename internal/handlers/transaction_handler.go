package handlers

import (
	"block-chat/internal/model"
	"block-chat/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func ValidateTransaction(c *gin.Context, myNode *model.Node) {
	var request model.Transaction

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sig_ok, err := myNode.Wallet.VerifySignature(request.Data, request.Signature, request.SenderAddress)

	if err != nil {
		log.Println("Error validated signature", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not verify transaction"})
		return
	}

	if sig_ok {
		log.Println("Signature was validated")
	}

	for _, node := range myNode.Ring {
		if node.PublicKey == request.SenderAddress {
			senderBalance := node.SoftBalance - node.SoftStake
			if senderBalance < request.CalculateFee() {
				log.Println("Insufficient funds to send transaction")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Insufficient funds to send transaction",
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction validated",
	})

	typeOfData, err := strconv.ParseBool(fmt.Sprint(request.TypeOfTransaction))

	if err != nil {
		log.Println(err)
	}

	receivedTransaction := model.Transaction{
		SenderAddress:     request.SenderAddress,
		ReceiverAddress:   request.ReceiverAddress,
		TypeOfTransaction: typeOfData,
		Data:              request.Data,
		Nonce:             myNode.Wallet.AddTransaction(),
		TransactionID:     "",
		Signature:         request.Signature,
	}
	//isBlockFull := myNode.CurrentBlock.AddTransaction(receivedTransaction, myNode)
	log.Println("Receives Txn")
	model.Mtx.Lock()

	txnPoolObj := model.TransactionInPool{Txn: receivedTransaction, IsSendTxn: false}
	model.TransactionPool <- txnPoolObj
	log.Println("Pushes txn to the txnPool")

	model.Mtx.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction received",
	})
	//if isBlockFull {
	//	myNode.CurrentBlock.ElectLeader(myNode)
	//}
	// Send response

}

// Send a transaction to another node
// ===================================
func SendTransaction(c *gin.Context, myNode *model.Node) {
	var request model.SendTransactionRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON. Error Message : " + err.Error())
		return
	}
	recipientId := request.Recipient
	recipientPublicAddress := utils.FindPublicAddress(myNode.Ring, recipientId)
	if recipientPublicAddress == nil {
		log.Println("Error : The given recipient id (=" + strconv.Itoa(recipientId) + ") does not correspond to any active node!")
	}

	typeOfData, err := strconv.ParseBool(fmt.Sprint(request.Message_or_Bitcoin))
	if err != nil {
		log.Println(err)
	}

	newTransaction := model.NewTransaction(
		recipientPublicAddress,
		typeOfData,
		request.Data,
		myNode.Wallet.AddTransaction(),
	)
	newTransaction.Signature, err = myNode.Wallet.SignTransaction(newTransaction)
	newTransaction.SenderAddress = myNode.Wallet.PublicKey
	if err != nil {
		log.Println("Error signing transaction", err)
	}

	log.Println("Sending transaction", newTransaction)

	transactionFee := newTransaction.CalculateFee()
	for _, node := range myNode.Ring {
		if node.Id == myNode.Id {
			senderBalance := node.SoftBalance - node.Stake
			if senderBalance < transactionFee {
				log.Println("Insufficient funds to send transaction")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Insufficient funds to send transaction",
				})
				return
			}
		}
	}
	model.Mtx.Lock()
	txnPoolObj := model.TransactionInPool{Txn: *newTransaction, IsSendTxn: true}
	model.TransactionPool <- txnPoolObj
	model.Mtx.Unlock()
	//isBlockFull := myNode.CurrentBlock.AddTransaction(*newTransaction, myNode)
	if myNode.BroadcastTransaction(newTransaction) {
		log.Println("Transaction broadcast successful.")
		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction sent",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Transaction not sent",
		})
	}

}
