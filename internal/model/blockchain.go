package model

import (
	"encoding/json"
	"fmt"
)

// Blockchain struct defines the blockchain with a slice of Blocks
type Blockchain struct {
	Chain []Block `json:"chain"`
}

// Blockchain toString()
func (bc *Blockchain) String() string {
	result := "Chain:\n"
	for ind, block := range bc.Chain {
		result += fmt.Sprintf("\n\t\t\t\t\t Block %d : %s",
			ind, block.String())

	}
	result += "\n"
	return result
}

// NewBlockchain creates a new Blockchain with the given initial blocks
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: nil,
	}
}

// JSONify serializes the Blockchain into a JSON string
func (bc *Blockchain) JSONify() (string, error) {
	jsonBytes, err := json.Marshal(bc)
	return string(jsonBytes), err
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) ValidateBlock(block *Block, node *Node) bool {
	for i := range block.Transactions {
		block.AddValidatedTransaction(block.Transactions[i], node)
	}
	node.SoftStateEqualToHardState()
	bc.Chain = append(bc.Chain, *block)
	newBlock := NewBlock(block.Index+1, block.CurrentHash)
	node.CurrentBlock = &newBlock
	return true
}

func (bc *Blockchain) AddNewBlock() Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	block := NewBlock(lastBlock.Index+1, lastBlock.CurrentHash)
	bc.Chain = append(bc.Chain, block)
	return bc.Chain[len(bc.Chain)-1]
}

// GetLastBlock return last block
func (bc *Blockchain) GetLastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}
