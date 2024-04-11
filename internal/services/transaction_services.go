package services

import (
	"block-chat/internal/model"
	"log"
)

func ValidateTxnService(myNode *model.Node, txn model.Transaction) {
	isBlockFull := myNode.CurrentBlock.AddTransaction(txn, myNode)
	if isBlockFull {
		log.Println("Block is full from validate and txn.Id : ", txn.TransactionID)
		model.IsBlockValidating = true
		myNode.CurrentBlock.ElectLeader(myNode)
	}
}
