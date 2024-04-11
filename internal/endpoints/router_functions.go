package endpoints

import (
	"block-chat/internal/model"
	"block-chat/internal/services"
	"log"
)

func processTransactions(myNode *model.Node) {
	for {
		if model.IsBlockValidating {
			log.Println("Waiting for the block validating signal!")
			<-model.BlockValidationSignal
			log.Println("Received the block validating signal!")
		}

		model.Mtx.Lock()
		for len(model.TransactionPool) > 0 {
			log.Println("Before len(model.TransactionPool) : ", len(model.TransactionPool))
			txnInPool := <-model.TransactionPool
			log.Println("After len(model.TransactionPool) : ", len(model.TransactionPool))
			services.ValidateTxnService(myNode, txnInPool.Txn)
			if model.IsBlockValidating == true {
				break
			}
		}
		model.Mtx.Unlock()
	}
}
