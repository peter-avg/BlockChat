package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Get Last Block from Blockchain
// ==============================
func GetLastBlock(c *gin.Context, MyNode *model.Node) {
	//var lastBlock = MyNode.Chain.GetLastBlock()
	var lastBlock = MyNode.CurrentBlock
	var responseString = "Last Block :\n\t" + lastBlock.String()
	c.String(http.StatusOK, responseString)
}

func ReceiveValidatedBlock(c *gin.Context, MyNode *model.Node) {
	// TODO : Update State of the Block
}
