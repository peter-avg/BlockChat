package handlers

import (
	"block-chat/internal/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Get Last Block from Blockchain
// ==============================
func GetLastBlock(c *gin.Context, MyNode *model.Node) {
	jsonBlock, err := json.Marshal(MyNode.Chain.GetLastBlock())
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"last_block": string(jsonBlock),
	})
}
