package handlers

import (
	"block-chat/internal/model"
	"block-chat/internal/utils"
	"crypto/rsa"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Is sent by a node to the bootstrap node to register itself, and its info gets broadcasted to all nodes
// ======================================================================================================
func RegisterNode(c *gin.Context, myNode *model.Node) {
	var request model.RegisterNodeRequest
	var wg sync.WaitGroup

	// Bind the request body to the request struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a public key from the modulus and exponent
	publicKey := rsa.PublicKey{
		N: request.Modulus,
		E: request.Exponent,
	}
	log.Println("myNode.Wallet = " + utils.Float64ToString(myNode.Wallet.Balance))
	if myNode.Wallet.Balance < 1000 {
		log.Println("Not enough money for the initial transaction")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not deduct money"})
		return
	}
	// Add the new node to the Ring
	NewNodeInfo := model.NewNodeInfo(myNode.Nonce+1, request.IP, request.Port, &publicKey, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		myNode.BroadcastNewNode(NewNodeInfo)
	}()
	newTransaction := model.NewTransaction(
		&publicKey,
		true,
		"1000",
		myNode.Wallet.AddTransaction(),
	)
	var err error
	newTransaction.Signature, err = myNode.Wallet.SignTransaction(newTransaction)
	newTransaction.SenderAddress = myNode.Wallet.PublicKey
	wg.Add(1)
	go func() {
		defer wg.Done()
		wg.Wait()
		if myNode.BroadcastTransaction(newTransaction) {
			log.Println("Initial Transaction broadcast successful.")
			//myNode.CurrentBlock.AddTransaction(*newTransaction, config.CAPACITY)
			c.JSON(http.StatusOK, gin.H{
				"message": "Transaction sent",
			})
			return
		}
	}()

	myNode.AddNewInfo(NewNodeInfo)
	myNode.CurrentBlock.AddTransaction(*newTransaction, myNode)

	// Serialize data for response
	jsonDataID := myNode.Nonce
	log.Println("jsonDataId : " + strconv.Itoa(jsonDataID))
	jsonBlockchain, err := json.Marshal(myNode.Chain)
	if err != nil {
		log.Println(err)
	}
	jsonRing, err := json.Marshal(myNode.Ring)
	if err != nil {
		log.Println(err)
	}

	//log.Println("myNode : ", myNode.String())

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"id":         jsonDataID,
		"blockchain": string(jsonBlockchain),
		"ring":       string(jsonRing),
		"balance":    1000,
	})

}

// Receive new node
// ================
func ReceiveNewNode(c *gin.Context, MyNode *model.Node) {
	var request model.NodeInfo

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	MyNode.AddNewInfo(&request)
	log.Println("Added new node to the Ring", MyNode.Ring)

	c.JSON(http.StatusOK, gin.H{
		"message": "Node added",
	})
}
