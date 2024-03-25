package handlers

import (
	"block-chat/internal/config"
	"block-chat/internal/model"
	"block-chat/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func ValidateTransaction(c *gin.Context, MyNode *model.Node) {
	var request model.Transaction

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sig_ok, err := MyNode.Wallet.VerifySignature(request.Data, request.Signature, request.SenderAddress)

	if err != nil {
		log.Println("Error validated signature", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not verify transaction"})
		return
	}

	if sig_ok {
		log.Println("Signature was validated")
	}

	for _, node := range MyNode.Ring {
		if node.PublicKey == request.SenderAddress {
			senderBalance := node.Balance - node.Stake
			if senderBalance < request.CalculateFee() {
				log.Println("Insufficient funds to send transaction")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Insufficient funds to send transaction",
				})
				return
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
			Nonce:             MyNode.Wallet.AddTransaction(),
			TransactionID:     "",
			Signature:         request.Signature,
		}

		MyNode.CurrentBlock.AddTransaction(receivedTransaction, config.CAPACITY)

		// Send response
		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction received",
		})
	}
}

// Send a transaction to another node
// ===================================
func SendTransaction(c *gin.Context, MyNode *model.Node) {
	var request model.SendTransactionRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON. Error Message : " + err.Error())
		return
	}
	recipientId := request.Recipient
	recipientPublicAddress := utils.FindPublicAddress(MyNode.Ring, recipientId)
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
		MyNode.Wallet.AddTransaction(),
	)
	newTransaction.Signature, err = MyNode.Wallet.SignTransaction(newTransaction)
	newTransaction.SenderAddress = MyNode.Wallet.PublicKey
	if err != nil {
		log.Println("Error signing transaction", err)
	}

	log.Println("Sending transaction", newTransaction)

	transactionFee := newTransaction.CalculateFee()
	if transactionFee > MyNode.Wallet.Balance {
		log.Println("Insufficient funds to send transaction")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient funds to send transaction",
		})
		return
	}

	if MyNode.BroadcastTransaction(newTransaction) {
		log.Println("Transaction broadcast successful.")
		MyNode.CurrentBlock.AddTransaction(*newTransaction, config.CAPACITY)
		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction sent",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Transaction not sent",
	})

}

//// Verify Transaction
//// ==================
//func ValidateTransaction(c *gin.Context, MyNode *model.Node) {
//	var request model.Transaction
//
//	if err := c.BindJSON(&request); err != nil {
//		log.Println("Error binding JSON")
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	sig_ok, err := MyNode.Wallet.VerifySignature(request.Data, request.Signature, request.SenderAddress)
//
//	if err != nil {
//		log.Println("Error validated signature", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not verify transaction"})
//		return
//	}
//
//	if sig_ok {
//		log.Println("Signature was validated")
//	}
//
//	for _, node := range MyNode.Ring {
//		if node.PublicKey == request.SenderAddress {
//			sender_balance := node.Balance - node.Stake
//			if sender_balance < request.CalculateFee() {
//				log.Println("Insufficient funds to send transaction")
//				c.JSON(http.StatusBadRequest, gin.H{
//					"error": "Insufficient funds to send transaction",
//				})
//				return
//			}
//
//			node.Balance -= request.CalculateFee()
//		}
//
//		c.JSON(http.StatusOK, gin.H{
//			"message": "Transaction validated",
//		})
//	}
//}

// Receive Transaction
// ===================
func ReceiveTransaction(c *gin.Context, MyNode *model.Node) {
	var request model.Transaction

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	typeOfData, err := strconv.ParseBool(fmt.Sprint(request.TypeOfTransaction))

	if err != nil {
		log.Println(err)
	}

	receivedTransaction := model.Transaction{
		SenderAddress:     request.SenderAddress,
		ReceiverAddress:   request.ReceiverAddress,
		TypeOfTransaction: typeOfData,
		Data:              request.Data,
		Nonce:             MyNode.Wallet.AddTransaction(),
		TransactionID:     "",
		Signature:         request.Signature,
	}

	MyNode.CurrentBlock.AddTransaction(receivedTransaction, config.CAPACITY)

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction received",
	})

}
