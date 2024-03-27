package services

//// Send a transaction to another node
//// ===================================
//func SendTransactionService(transaction *model.Transaction, myNode *model.Node)  {
//	var err error
//	transaction.Signature, err = myNode.Wallet.SignTransaction(transaction)
//	transaction.SenderAddress = myNode.Wallet.PublicKey
//	if err != nil {
//		log.Println("Error signing transaction", err)
//	}
//
//	log.Println("Sending transaction", transaction)
//
//	transactionFee := transaction.CalculateFee()
//	if transactionFee > myNode.Wallet.Balance {
//		log.Println("Insufficient funds to send transaction")
//		utils.CustomError{}
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": "Insufficient funds to send transaction",
//		})
//		return
//	}
//	myNode.CurrentBlock.AddTransaction(*transaction, myNode)
//	if myNode.BroadcastTransaction(transaction) {
//		log.Println("Transaction broadcast successful.")
//		//myNode.CurrentBlock.AddTransaction(*newTransaction, config.CAPACITY)
//		c.JSON(http.StatusOK, gin.H{
//			"message": "Transaction sent",
//		})
//		return
//	}
//
//	c.JSON(http.StatusBadRequest, gin.H{
//		"error": "Transaction not sent",
//	})
//
//}
