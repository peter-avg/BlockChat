package blockchain

import (
	// "crypto/rsa"
	// "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	// "crypto/x509"
	"encoding/hex"
	"encoding/json"
	// "encoding/pem"
	// "fmt"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	SenderAddress      *rsa.PublicKey    `json:"sender_address"`
	ReceiverAddress    int    `json:"receiver_address"`
	TypeOfTransaction  bool   `json:"type_of_transaction"` // 0 for message, 1 for bcc
	Data               string `json:"data"`
	Nonce              int    `json:"nonce"`
	TransactionID      string  `json:"transaction_id"`
	Signature          []byte `json:"signature"`
}

// NewTransaction creates a new Transaction
func NewTransaction(receiverAddress int, typeOfTransaction bool, data string, nonce int) *Transaction {
    
	return &Transaction{
		ReceiverAddress:    receiverAddress,
		TypeOfTransaction:  typeOfTransaction,
		Data:               data,
		Nonce:              nonce,
		TransactionID:      "",
		Signature:          nil,
	}
}

// JSONify serializes the Transaction into a JSON string
func (t *Transaction) JSONify() string {
	jsonBytes, err := json.Marshal(t)
    if err != nil {
        println("Could not jsonify transaction");
        return ""
    }
	return string(jsonBytes)
}

// Hashify creates a hash for the Transaction object
func (t *Transaction) Hashify() {
	jsonString := t.JSONify()
	hash := sha256.Sum256([]byte(jsonString))
	t.TransactionID = hex.EncodeToString(hash[:]);
}

