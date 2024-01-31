package blockchain

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "math/big"
    "crypto/rsa"
    "fmt"
)


// Is sent by a node to the bootstrap node to register itself
// The node receives back the blockchain, the node's id in the Ring, and its balance
func RegisterNode(c *gin.Context) {
    var request RegisterNodeRequest;

    // Bind the request body to the request struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // Create a public key from the modulus and exponent
	NBigInt := new(big.Int)
	NBigInt.SetString(request.Modulus, 10)
	publicKey := rsa.PublicKey{
		N: NBigInt,
		E: request.Exponent,
	}

    // Add the new node to the Ring
    NewNodeInfo := NewNodeInfo(MyNode.Ring[len(MyNode.Ring)-1].Id + 1, request.IP, request.Port, &publicKey, 1000);
    MyNode.AddNewInfo(NewNodeInfo);

    responseData := gin.H{
        "ring" : MyNode.Ring,
		"blockchain": MyNode.Chain,
		"balance":    1000,
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        fmt.Sprint(MyNode.Ring[len(MyNode.Ring)-1].Id),
		"blockchain": MyNode.Chain,
		"balance":    1000,
	})
};

