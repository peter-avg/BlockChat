package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Block represents each 'item' in the blockchain
type Block struct {
    Index         int           `json:"index"`
    Timestamp     int64         `json:"timestamp"`
    Transactions  []Transaction `json:"transactions"`
    Validator     int           `json:"validator"`
    CurrentHash   string        `json:"current_hash"`
    PreviousHash  string        `json:"previous_hash"`
}

// NewBlock creates and returns a new Block
func NewBlock(index int, previousHash string) *Block {
    t := time.Now();
    timestamp := t.UnixNano();
    return &Block{
        Index:         index,
        Timestamp:     timestamp,
        Transactions:  nil,
        Validator:     0,
        PreviousHash:  previousHash,
        CurrentHash:   "",
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

// AddTransaction adds a new transaction to the block if there's capacity
func (b *Block) AddTransaction(transaction *Transaction, capacity int) {
    if len(b.Transactions) < capacity {
        b.Transactions = append(b.Transactions, *transaction)
    }
}

// Set previous Hash of a block
func (b *Block) GetPreviousHash(chain *Blockchain) {
    b.PreviousHash = chain.Chain[len(chain.Chain) - 1].CurrentHash;
}

