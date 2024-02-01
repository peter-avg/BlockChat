package main

import (
	// "fmt"
	"api/blockchain/blockchain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	// "math/big"
	"crypto/rsa"
)

// Is sent by a node to the bootstrap node to register itself
// The node receives back the blockchain, the node's id in the Ring, and its balance
// =================================================================================
func RegisterNode(c *gin.Context) {
    var request blockchain.RegisterNodeRequest;

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

    // Add the new node to the Ring
    NewNodeInfo := blockchain.NewNodeInfo(MyNode.Nonce, request.IP, request.Port, &publicKey, 1000);
    MyNode.AddNewInfo(NewNodeInfo);
    log.Println("Added new node to the Ring", MyNode.Ring);

    // Serialize data for response
    jsonDataID := MyNode.Nonce;
    jsonBlockchain,err := json.Marshal(MyNode.Chain);
    if err != nil {
        log.Println(err);
    }
    jsonRing,err := json.Marshal(MyNode.Ring);
    if err != nil {
        log.Println(err);
    }

    MyNode.Wallet.DeductMoney(1000);

    // Send response
	c.JSON(http.StatusOK, gin.H{
		"id": jsonDataID,
		"blockchain": string(jsonBlockchain),
        "ring" : string(jsonRing),
        "balance":    1000,
	})
};

// Send a transaction to another node
// ===================================
func SendTransaction(c *gin.Context) {
    // var request blockchain.SendTransactionRequest;

    // Send response
    c.JSON(http.StatusOK, gin.H{
        "message": "Transaction sent",
    })
}

// Set Stake for a block
// ===================================
func SetStake(c *gin.Context) {
    // var request blockchain.SetStakeRequest;

    // Send response
    c.JSON(http.StatusOK, gin.H{
        "message": "Stake set",
    })
}

// Get Balance from Wallet
// ===================================
func GetBalance(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "balance": fmt.Sprint(MyNode.Wallet.Balance),
    })
}

// Get Last Block from Blockchain
// ===================================
func GetLastBlock(c *gin.Context) {
    jsonBlock,err := json.Marshal(MyNode.Chain.GetLastBlock());
    if err != nil {
        log.Println(err);
    }
    c.JSON(http.StatusOK, gin.H{
        "last_block": string(jsonBlock),
    })
}

