package model

import (
	"block-chat/internal/config"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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
	return fmt.Sprintf("Index: %d, Timestamp: %s,\n\tTransactions: %s,\n\tValidator: %d, CurrentHash: %s, PreviousHash: %s",
		b.Index, timeString, transactionsString, b.Validator, b.CurrentHash, b.PreviousHash)
}

// AddTransaction adds a new transaction to the block if there's capacity
func (b *Block) AddTransaction(transaction Transaction, myNode *Node) {
	var capacity = config.CAPACITY
	var numberOfTransactions = len(b.Transactions)
	var transactionFee = transaction.CalculateFee()
	var isStakeTransaction = false
	if transaction.ReceiverAddress.Equal(config.STAKE_PUBLIC_ADDRESS) {
		isStakeTransaction = true
	}
	// TODO: refresh myNode.wallet & nodeInfo.balance
	// TODO: add the transaction to the current block
	if myNode.Wallet.PublicKey.Equal(transaction.ReceiverAddress) {
		if transaction.TypeOfTransaction == true {
			myNode.Wallet.Balance += transactionFee
			log.Println("\t\tBlock Chat Coins Received!\n\t\t--------------------\nYou got transferred " + strconv.FormatFloat(transaction.CalculateFee(), 'f', -1, 64) + " Block Chat Coins!")
		} else {
			log.Println("\t\tMessage Received!\n\t\t--------------------\n" + transaction.Data)
		}
	}

	if myNode.Wallet.PublicKey.Equal(transaction.SenderAddress) {
		myNode.Wallet.Balance -= transactionFee
		if transaction.TypeOfTransaction == true {
			log.Println("\t\tBlock Chat Coins Sent!\n\t\t--------------------\nYou sent " + strconv.FormatFloat(transaction.CalculateFee(), 'f', -1, 64) + " Block Chat Coins!")
		} else {
			log.Println("\t\tMessage Sent!\n\t\t--------------------\n" + transaction.Data)
		}
	}

	for _, nodeInfo := range myNode.Ring {
		if nodeInfo.PublicKey.Equal(transaction.ReceiverAddress) && transaction.TypeOfTransaction == true {
			nodeInfo.Balance += transactionFee
		}
		if nodeInfo.PublicKey.Equal(transaction.SenderAddress) {
			nodeInfo.Balance -= transactionFee
			if isStakeTransaction {
				nodeInfo.Stake += transactionFee
			}
		}
	}

	// TODO: process if block is not full
	if numberOfTransactions+1 < capacity {
		numberOfTransactions++
		return
	}

	// TODO : process if block is full
	// TODO : elect new leader
	// TODO : implement new endpoint for the leader to validate block
	// TODO : implement new endpoint for other nodes to receive
	// 		  & add the validated block onto their blockchain
}
