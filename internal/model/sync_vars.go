package model

import "sync"

type TransactionInPool struct {
	Txn       Transaction
	IsSendTxn bool
	// if valOrSend == false its ValidateTxn
	// if valOrSend == true its SendTxn
}

// GLOBAL SYNCHRONIZATION VARIABLES
var TransactionPool = make(chan TransactionInPool, 100)

var IsBlockValidating = false

var BlockValidationSignal = make(chan struct{}, 1)

var Mtx sync.Mutex
