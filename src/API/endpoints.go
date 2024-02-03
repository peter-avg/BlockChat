package main

import (
	// "fmt"
	"api/blockchain/blockchain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
    var request blockchain.SendTransactionRequest;

    if err := c.BindJSON(&request); err != nil {
        log.Println("Error binding JSON");
    }

    receiver,err := strconv.Atoi(string(request.Recipient[len(request.Recipient)-1]));
    if err != nil {
        log.Println(err);
    }

    type_of_data,err := strconv.ParseBool(fmt.Sprint(request.Message_or_Bitcoin));
    if err != nil { 
        log.Println(err);
    }

    new_transaction := blockchain.NewTransaction(
        receiver,
        type_of_data,
        request.Data,
        MyNode.Wallet.AddTransaction(),
    );
    new_transaction.Signature, err = MyNode.Wallet.SignTransaction(new_transaction);
    new_transaction.SenderAddress = MyNode.Wallet.PublicKey;
    if err != nil {
        log.Println("Error signing transaction", err);
    }

    log.Println("Sending transaction", new_transaction);

    MyNode.BroadcastTransaction(new_transaction);

    // Send response
    c.JSON(http.StatusOK, gin.H{
        "message": "Transaction sent",
    })
}

// Set Stake for Proof of Stake
// ============================
func SetStake(c *gin.Context) {
    var request blockchain.SetStakeRequest;

    if err := c.BindJSON(&request); err != nil {
        log.Println("Error binding JSON");
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    MyNode.Stake = request.Stake;
    
    // Send response
    c.JSON(http.StatusOK, gin.H{
        "message": "Stake set",
    })
}

// Get Balance from Wallet
// =======================
func GetBalance(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "balance": fmt.Sprint(MyNode.Wallet.Balance),
    })
}

// Get Last Block from Blockchain
// ==============================
func GetLastBlock(c *gin.Context) {
    jsonBlock,err := json.Marshal(MyNode.Chain.GetLastBlock());
    if err != nil {
        log.Println(err);
    }
    c.JSON(http.StatusOK, gin.H{
        "last_block": string(jsonBlock),
    })
}

// Receive Transaction
// ===================
func ReceiveTransaction(c *gin.Context) {
    var request blockchain.Transaction;

    if err := c.BindJSON(&request); err != nil {
        log.Println("Error binding JSON");
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    type_of_data,err := strconv.ParseBool(fmt.Sprint(request.TypeOfTransaction));

    if err != nil { 
        log.Println(err);
    }

    received_transaction := blockchain.NewTransaction(
        MyNode.Id,
        type_of_data,
        request.Data,
        MyNode.Wallet.Nonce + 1,
    );
    received_transaction.Signature = request.Signature;
    received_transaction.Data = request.Data;
    received_transaction.SenderAddress = request.SenderAddress;
    log.Println("Received transaction", received_transaction);

    verification,err := MyNode.Wallet.VerifyTransaction(received_transaction);
    if err != nil {
        log.Println("Error verifying transaction", err);
        return
    }

    MyNode.CurrentBlock.AddTransaction(received_transaction,CAPACITY);

    if verification { 
        log.Println("Transaction verified");
    }

    log.Println("Received transaction", received_transaction);

    // Send response
    c.JSON(http.StatusOK, gin.H{
        "message": "Transaction received",
    })

}


