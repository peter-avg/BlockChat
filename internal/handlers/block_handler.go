package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Get Last Block from Blockchain
// ==============================
func GetLastBlock(c *gin.Context, MyNode *model.Node) {
	var lastValidatedBlock = MyNode.Chain.GetLastBlock()
	var responseString = ""
	responseString += "Last Validated Block :\n\t" + lastValidatedBlock.String() + "\n"
	var currentBlock = MyNode.CurrentBlock
	responseString = "Current Block :\n\t" + currentBlock.String()
	c.String(http.StatusOK, responseString)
}

func ReceiveValidatedBlock(c *gin.Context, myNode *model.Node) {
	var request model.Block

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myNode.Chain.ValidateBlock(&request, myNode)

	c.JSON(http.StatusOK, gin.H{
		"message": "Block Validated",
	})
}
