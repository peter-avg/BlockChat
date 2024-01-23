package blockchain

import (
    "encoding/json"
)

// Blockchain struct defines the blockchain with a slice of Blocks
type Blockchain struct {
    Chain []Block `json:"chain"`
}

// NewBlockchain creates a new Blockchain with the given initial blocks
func NewBlockchain(initialBlocks []Block) *Blockchain {
    return &Blockchain{
        Chain: initialBlocks,
    }
}

// JSONify serializes the Blockchain into a JSON string
func (bc *Blockchain) JSONify() (string, error) {
    jsonBytes, err := json.Marshal(bc)
    return string(jsonBytes), err
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block Block) {
    bc.Chain = append(bc.Chain, block)
}

