package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// Get Last Block from Blockchain
// ==============================
func GetLastBlock(c *gin.Context, MyNode *model.Node) {
	//var lastValidatedBlock = MyNode.Chain.GetLastBlock()
	var lastValidatedBlock = MyNode.CurrentBlock
	var responseString = "Last Block :\n\t" + lastValidatedBlock.String()
	c.String(http.StatusOK, responseString)
}

func ReceiveValidatedBlock(c *gin.Context, myNode *model.Node) {
	var request model.Block

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Received Block Id : " + strconv.Itoa(request.Index))
	myNode.Chain.ValidateBlock(&request, myNode)
	log.Println("New Block Added : " + request.String())

	c.JSON(http.StatusOK, gin.H{
		"message": "Block Validated",
	})
}
