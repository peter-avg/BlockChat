package model

import (
	"block-chat/internal/config"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"math/rand"
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
func (b Block) String() string {
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
	log.Println("Current Block Size : " + strconv.Itoa(numberOfTransactions))
	var transactionFee = transaction.CalculateFee()
	log.Println("Transaction Fee : " + strconv.FormatFloat(transactionFee, 'f', -1, 64))
	var isStakeTransaction = false
	if transaction.ReceiverAddress.Equal(config.STAKE_PUBLIC_ADDRESS) {
		isStakeTransaction = true
	}
	var isUnstakeTransaction = false
	if isStakeTransaction && transactionFee == 0 {
		isUnstakeTransaction = true
		isStakeTransaction = false
	}

	// TODO: refresh myNode.wallet & nodeInfo.balance
	// TODO: add the transaction to the current block
	if myNode.Wallet.PublicKey.Equal(transaction.ReceiverAddress) {
		if transaction.TypeOfTransaction == true {
			//myNode.Wallet.Balance += transactionFee
			log.Println("\t\tBlock Chat Coins Received!\n\t\t--------------------\n\t\tYou got transferred " + strconv.FormatFloat(transaction.CalculateFee(), 'f', -1, 64) + " Block Chat Coins!")
		} else {
			log.Println("\t\t\t\tMessage Received!\n\t\t--------------------\n" + transaction.Data)
		}
	}

	if myNode.Wallet.PublicKey.Equal(transaction.SenderAddress) {
		//myNode.Wallet.Balance -= transactionFee
		if transaction.TypeOfTransaction == true {
			log.Println("\t\tBlock Chat Coins Sent!\n\t\t--------------------\n\t\tYou sent " + strconv.FormatFloat(transaction.CalculateFee(), 'f', -1, 64) + " Block Chat Coins!")
		} else {
			log.Println("\t\tMessage Sent!\n\t\t--------------------\n\t\tMessage : " + transaction.Data)
		}
	}

	for i, nodeInfo := range myNode.Ring {
		if nodeInfo.PublicKey.Equal(transaction.ReceiverAddress) && transaction.TypeOfTransaction == true {
			myNode.Ring[i].SoftBalance += transactionFee
		}
		if nodeInfo.PublicKey.Equal(transaction.SenderAddress) {
			myNode.Ring[i].SoftBalance -= transactionFee
			if isStakeTransaction {
				myNode.Ring[i].SoftStake += transactionFee
			}
			if isUnstakeTransaction {
				myNode.Ring[i].SoftBalance += myNode.Ring[i].SoftStake
				myNode.Ring[i].SoftStake = 0
			}
		}
	}

	// TODO: process if block is not full
	if numberOfTransactions+1 < capacity {
		b.Transactions = append(b.Transactions, transaction)
		return
	}

	// TODO : process if block is full
	// TODO : maybe define mintBlock()
	// TODO : elect new leader
	var leaderId = 0
	var totalStakeAmount float64 = 0
	for _, nodeInfo := range myNode.Ring {
		totalStakeAmount += nodeInfo.Stake
	}
	log.Printf("Total Stake Amount : %f\n", totalStakeAmount)
	if totalStakeAmount != 0 {
		var seedString = myNode.Chain.GetLastBlock().CurrentHash
		hash := fnv.New64()
		hash.Write([]byte(seedString))
		seed := int64(hash.Sum64())
		rand.Seed(seed)
		roundedNumberInt := int(math.Round(totalStakeAmount))
		randomNumber := rand.Intn(roundedNumberInt + 1)
		log.Printf("Random Generated Number : %d\n", randomNumber)
		var currSum float64 = 0
		for _, nodeInfo := range myNode.Ring {
			currSum += nodeInfo.Stake
			if int(math.Round(currSum)) >= randomNumber {
				leaderId = nodeInfo.Id
				break
			}
		}
	}
	log.Printf("Elected Leader Node Id : %d\n", leaderId)

	// TODO : implement new endpoint for the leader to validate block
	// TODO : implement new endpoint for other nodes to receive
	// 		  & add the validated block onto their blockchain

	if myNode.Id == leaderId {
		log.Printf("I am the leader node with myNode.id == %d\n", myNode.Id)
		b.MintBlock(myNode)
	}
}

func (b *Block) MintBlock(myNode *Node) {
}

// AddTransaction adds a new transaction to the block if there's capacity
func (b *Block) AddValidatedTransaction(transaction Transaction, myNode *Node) {
	var transactionFee = transaction.CalculateFee()
	log.Println("Transaction Fee : " + strconv.FormatFloat(transactionFee, 'f', -1, 64))
	var isStakeTransaction = false
	if transaction.ReceiverAddress.Equal(config.STAKE_PUBLIC_ADDRESS) {
		isStakeTransaction = true
	}
	var isUnStakeTransaction = false
	if isStakeTransaction && transactionFee == 0 {
		isUnStakeTransaction = true
		isStakeTransaction = false
	}
	if myNode.Wallet.PublicKey.Equal(transaction.ReceiverAddress) {
		if transaction.TypeOfTransaction == true {
			myNode.Wallet.Balance += transactionFee
		}
	}

	if myNode.Wallet.PublicKey.Equal(transaction.SenderAddress) {
		myNode.Wallet.Balance -= transactionFee
	}

	for i, nodeInfo := range myNode.Ring {
		if nodeInfo.PublicKey.Equal(transaction.ReceiverAddress) && transaction.TypeOfTransaction == true {
			myNode.Ring[i].Balance += transactionFee
		}
		if nodeInfo.PublicKey.Equal(transaction.SenderAddress) {
			myNode.Ring[i].Balance -= transactionFee
			if isStakeTransaction {
				myNode.Ring[i].Stake += transactionFee
			}
			if isUnStakeTransaction {
				myNode.Ring[i].Balance += myNode.Ring[i].Stake
				myNode.Ring[i].Stake = 0
			}
		}
	}
}
