package services

//// Set Stake for Proof of Stake
//// ============================
//func SendStake(c *gin.Context, MyNode *model.Node) {
//	// Stake is a Transaction with a Recipient Address == -1
//
//	var request model.SetStakeRequest
//
//	if err := c.BindJSON(&request); err != nil {
//		log.Println("Error binding JSON")
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if MyNode.Wallet.DeductMoney(request.Stake) == false {
//		log.Println("Could not set stake, insufficient funds")
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not set stake, insufficient funds"})
//		return
//	}
//
//	nodeInfo := model.NodeInfo{
//		Id:    MyNode.Id,
//		Stake: request.Stake,
//	}
//
//	if MyNode.BroadcastStake(nodeInfo) {
//		log.Println("Stake broadcasted")
//		log.Println("Setting stake", request.Stake)
//		MyNode.Ring[MyNode.Id].Stake = request.Stake
//		c.JSON(http.StatusOK, gin.H{
//			"message": "Stake set",
//		})
//		return
//	}
//
//	log.Println("Stake not set")
//	c.JSON(http.StatusBadRequest, gin.H{
//		"error": "Stake not set",
//	})
//}
//
//func SetStake(c *gin.Context, MyNode *model.Node) {
//	// SendStake() creates and broadcasts a Transaction
//	// where : receiverAddress is equal to -1
//	// 		   typeOfTransaction is equal to True (BCC)
//	//
//	var request model.SetStakeRequest
//	var err error
//	if err = c.BindJSON(&request); err != nil {
//		log.Println("Error binding JSON")
//	}
//
//	var stakeAmount int = request.Stake
//	if stakeAmount > MyNode.Wallet.Balance {
//		log.Println("Insufficient funds to send transaction")
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": "Insufficient funds to send transaction",
//		})
//		return
//	}
//
//	var newStakeTransaction *model.Transaction = model.NewTransaction(
//		-1,
//		true,
//		strconv.Itoa(stakeAmount),
//		MyNode.Wallet.AddTransaction(),
//	)
//
//	newStakeTransaction.Signature, err = MyNode.Wallet.SignTransaction(newStakeTransaction)
//	newStakeTransaction.SenderAddress = MyNode.Wallet.PublicKey
//	if err != nil {
//		log.Println("Error signing stake transaction", err)
//	}
//
//	if MyNode.BroadcastTransaction(newStakeTransaction) {
//		log.Println("Stake broadcasted")
//		MyNode.CurrentBlock.AddTransaction(*newStakeTransaction, config.CAPACITY)
//		c.JSON(http.StatusOK, gin.H{
//			"message": fmt.Sprint("Stake Transaction of ", stakeAmount, " is send")
//		})
//		return
//	}
//
//	c.JSON(http.StatusBadRequest, gin.H{
//		"error": "Stake Transaction not sent",
//	})
//
//}
//
//// Validate Stake
//// ==============
//func validateStake(c *gin.Context, MyNode *model.Node) {
//	var request model.NodeInfo
//
//	if err := c.BindJSON(&request); err != nil {
//		log.Println("Error binding JSON")
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if request.Stake > MyNode.Ring[request.Id].Balance {
//		log.Println("Insufficient funds to set stake")
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": "Insufficient funds to set stake",
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "Stake validated",
//	})
//}
//
//func ValidateStake(c *gin.Context, MyNode *model.Node) {
//	var requestTransaction model.Transaction
//
//	if err := c.BindJSON(&requestTransaction); err != nil {
//		log.Println("Error binding JSON")
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	sig_ok, err := MyNode.Wallet.VerifyTransaction(requestTransaction.Data,
//		requestTransaction.Signature,
//		requestTransaction.SenderAddress)
//
//	if err != nil {
//		log.Println("Error validating stake signature", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not verify stake transaction"})
//		return
//	}
//	if sig_ok {
//		log.Println("Stake signature was validated")
//	}
//
//	if request.Stake > MyNode.Ring[request.Id].Balance {
//		log.Println("Insufficient funds to set stake")
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": "Insufficient funds to set stake",
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "Stake validated",
//	})
//}
//
//// Receive Stake
//// =============
//func ReceiveStake(c *gin.Context, MyNode *model.Node) {
//	var request model.NodeInfo
//
//	if err := c.BindJSON(&request); err != nil {
//		log.Println("Error binding JSON")
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	MyNode.Ring[request.Id].Stake = request.Stake
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "Stake received",
//	})
//}
