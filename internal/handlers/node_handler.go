package handlers

import (
	"block-chat/internal/model"
	"crypto/rsa"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Is sent by a node to the bootstrap node to register itself, and its info gets broadcasted to all nodes
// ======================================================================================================
func RegisterNode(c *gin.Context, MyNode *model.Node) {
	var request model.RegisterNodeRequest

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

	if MyNode.Wallet.DeductMoney(1000) == false {
		log.Println("Error deducting money")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not deduct money"})
		return
	}

	// Add the new node to the Ring
	NewNodeInfo := model.NewNodeInfo(MyNode.Nonce, request.IP, request.Port, &publicKey, 1000)
	MyNode.AddNewInfo(NewNodeInfo)
	log.Println("Added new node to the Ring", MyNode.Ring)

	// Serialize data for response
	jsonDataID := MyNode.Nonce
	jsonBlockchain, err := json.Marshal(MyNode.Chain)
	if err != nil {
		log.Println(err)
	}
	jsonRing, err := json.Marshal(MyNode.Ring)
	if err != nil {
		log.Println(err)
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"id":         jsonDataID,
		"blockchain": string(jsonBlockchain),
		"ring":       string(jsonRing),
		"balance":    1000,
	})

	MyNode.BroadcastNewNode(NewNodeInfo)
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
