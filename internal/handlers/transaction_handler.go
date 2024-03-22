package handlers

import (
	"block-chat/internal/config"
	"block-chat/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// Send a transaction to another node
// ===================================
func SendTransaction(c *gin.Context, MyNode *model.Node) {
	var request model.SendTransactionRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON. Error Message : " + err.Error())
		return
	}
	receiver, err := strconv.Atoi(string(request.Recipient[len(request.Recipient)-1]))

	if err != nil {
		log.Println("request.Recipient : " + request.Recipient)
		log.Println(err)
	}

	type_of_data, err := strconv.ParseBool(fmt.Sprint(request.Message_or_Bitcoin))
	if err != nil {
		log.Println(err)
	}

	new_transaction := model.NewTransaction(
		receiver,
		type_of_data,
		request.Data,
		MyNode.Wallet.AddTransaction(),
	)
	new_transaction.Signature, err = MyNode.Wallet.SignTransaction(new_transaction)
	new_transaction.SenderAddress = MyNode.Wallet.PublicKey
	if err != nil {
		log.Println("Error signing transaction", err)
	}

	log.Println("Sending transaction", new_transaction)

	transaction_fee := new_transaction.CalculateFee()
	//if transaction_fee > MyNode.Wallet.Balance {
	if MyNode.Wallet.DeductMoney(transaction_fee) == false {
		log.Println("Insufficient funds to send transaction")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient funds to send transaction",
		})
		return
	}

	if MyNode.BroadcastTransaction(new_transaction) {
		log.Println("Transaction broadcasted")
		MyNode.CurrentBlock.AddTransaction(*new_transaction, config.CAPACITY)
		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction sent",
		})

		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Transaction not sent",
	})
}

// Verify Transaction
// ==================
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
			sender_balance := node.Balance - node.Stake
			if sender_balance < request.CalculateFee() {
				log.Println("Insufficient funds to send transaction")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Insufficient funds to send transaction",
				})
				return
			}

			node.Balance -= request.CalculateFee()
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction validated",
		})
	}
}

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
		ReceiverAddress:   MyNode.Id,
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
