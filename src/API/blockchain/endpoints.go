package blockchain

import (
	// "fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	// "math/big"
	// "crypto/rsa"
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
	// publicKey := rsa.PublicKey{
	// 	N: request.Modulus,
	// 	E: request.Exponent,
	// }

    // Add the new node to the Ring
    // NewNodeInfo := NewNodeInfo(MyNode.Ring[len(MyNode.Ring)-1].Id + 1, request.IP, request.Port, &publicKey, 1000);
    // MyNode.AddNewInfo(NewNodeInfo);

     //    responseData := gin.H{
     //        "ring" : MyNode.Ring,
	// 	"blockchain": MyNode.Chain,
	// 	"balance":    1000,
	// }

    jsonDataChain,_ := MyNode.Chain.JSONify()

	c.JSON(http.StatusOK, gin.H{
		"id":       string(MyNode.Ring[len(MyNode.Ring)-1].Id + 1),
		"blockchain": jsonDataChain,
        "balance":    1000,
	})
};

