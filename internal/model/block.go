package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Validator    int           `json:"validator"`
	CurrentHash  string        `json:"current_hash"`
	PreviousHash string        `json:"previous_hash"`
}

// NewBlock creates and returns a new Block
func NewBlock(index int, previousHash string) *Block {
	t := time.Now()
	timestamp := t.UnixNano()
	return &Block{
		Index:        index,
		Timestamp:    timestamp,
		Transactions: nil,
		Validator:    0,
		PreviousHash: previousHash,
		CurrentHash:  "",
	}
}

// JSONify serializes the Block into a JSON string
func (b *Block) JSONify() (string, error) {
	jsonBytes, err := json.Marshal(b)
	return string(jsonBytes), err
}

// Hashify updates the current hash of the block
func (b *Block) Hashify() {
	jsonString, err := b.JSONify()
	if err != nil {
		fmt.Println("Could not create hash:", err)
		return
	}
	hash := sha256.Sum256([]byte(jsonString))
	b.CurrentHash = hex.EncodeToString(hash[:])
}

// GetPreviousHash Set previous Hash of a block
func (b *Block) GetPreviousHash(chain *Blockchain) {
	b.PreviousHash = chain.Chain[len(chain.Chain)-1].CurrentHash
}

// GetTimestamp Get a timestamp
func GetTimestamp() int64 {
	t := time.Now()
	return t.UnixNano()
}

// Block toString()
func (b *Block) String() string {
	t := time.Unix(0, b.Timestamp)
	var timeString string = t.Format("15:04:05")
	transactionsString := ""
	for ind, transaction := range b.Transactions {
		transactionsString += fmt.Sprintf("\n\t\t Transaction %d : %s",
			ind, transaction.String())
	}
	return fmt.Sprintf("Index: %d, Timestamp: %s,\n\tTransactions: %s, Validator: %d, CurrentHash: %s, PreviousHash: %s",
		b.Index, timeString, transactionsString, b.Validator, b.CurrentHash, b.PreviousHash)
}

// AddTransaction adds a new transaction to the block if there's capacity
func (b *Block) AddTransaction(transaction Transaction, capacity int) {
	if len(b.Transactions) < capacity {
		b.Transactions = append(b.Transactions, transaction)
		return
	}

	fmt.Println("Block is full, cannot add transaction")
	// TODO: start proof of stake process
	// TODO: after proof of stake, get the new blockchain probably
	//blockchain.getLastBlock()
	// find the stakes of all the nodes in the block
	// from all the transactions with receiver_addr == 0
	var blockTransactions []Transaction = b.Transactions

	for ind, transaction := range blockTransactions {
		fmt.Println("transaction: ", ind, " : ", transaction)
	}

	// stake transactions are transactions with
	// receiver_address equal to 0
	var totalAmountStaked float64 = 0
	for _, transaction := range blockTransactions {
		if transaction.ReceiverAddress == -1 {
			var stakeAmount = transaction.CalculateFee()
			totalAmountStaked += stakeAmount
		}
	}
	// make the lottery
	// return the stakes to all(?) nodes
	// create new (empty)? block
	// and make it the end of the blockchain
	// TODO: create new block and add transaction there
}
