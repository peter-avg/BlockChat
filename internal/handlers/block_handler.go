package handlers

import (
	"block-chat/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Get Last Block from Blockchain
// ==============================
func GetLastBlock(c *gin.Context, MyNode *model.Node) {
	//log.Println("Before json.Marshal().")
	//jsonBlock, err := json.Marshal(MyNode.Chain.GetLastBlock())
	//log.Println("After json.Marshal().")
	//if err != nil {
	//	log.Println(err)
	//}
	var lastBlock model.Block = MyNode.Chain.GetLastBlock()
	var responseString string = "Last Block :\n\t" + lastBlock.String()
	c.String(http.StatusOK, responseString)

	//c.JSON(http.StatusOK, gin.H{
	//	//"last_block": string(jsonBlock),
	//	"last_block": lastBlock.String(),
	//})
}
