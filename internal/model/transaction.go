package model

import (
	"fmt"
	// "crypto/rsa"
	"strconv"
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
	SenderAddress     *rsa.PublicKey `json:"sender_address"`
	ReceiverAddress   *rsa.PublicKey `json:"receiver_address"`
	TypeOfTransaction bool           `json:"type_of_transaction"` // 0 for message, 1 for bcc
	Data              string         `json:"data"`
	Nonce             int            `json:"nonce"`
	TransactionID     string         `json:"transaction_id"`
	Signature         []byte         `json:"signature"`
}

// NewTransaction creates a new Transaction
func NewTransaction(receiverAddress *rsa.PublicKey, typeOfTransaction bool, data string, nonce int) *Transaction {
	return &Transaction{
		ReceiverAddress:   receiverAddress,
		TypeOfTransaction: typeOfTransaction,
		Data:              data,
		Nonce:             nonce,
		TransactionID:     "",
		Signature:         nil,
	}
}

func (t *Transaction) String() string {
	//return fmt.Sprintf("Transaction[SenderAddress: %v, ReceiverAddress: %d, TypeOfTransaction: %t, Data: %s, Nonce: %d, TransactionID: %s, Signature: %v]",
	//	t.SenderAddress, t.ReceiverAddress, t.TypeOfTransaction, t.Data, t.Nonce, t.TransactionID, t.Signature)
	return fmt.Sprintf("Transaction[SenderAddress: rsa_key,"+
		" ReceiverAddress: %v,"+
		" TypeOfTransaction: %t,"+
		" Data: %s,"+
		" Nonce: %d,"+
		" TransactionID: %s,"+
		" Signature: %v]",
		"t.ReceiverAddress", t.TypeOfTransaction, t.Data,
		t.Nonce, t.TransactionID, "signature_number")
}

// JSONify serializes the Transaction into a JSON string
func (t *Transaction) JSONify() string {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		println("Could not jsonify transaction")
		return ""
	}
	return string(jsonBytes)
}

// Hashify creates a hash for the Transaction object
func (t *Transaction) Hashify() {
	jsonString := t.JSONify()
	hash := sha256.Sum256([]byte(jsonString))
	t.TransactionID = hex.EncodeToString(hash[:])
}

// CalculateFee Calculate transaction fee for messages of bitcoin
func (t *Transaction) CalculateFee() float64 {
	// if it's bitcoin, the fee is the data
	if t.TypeOfTransaction {
		res, err := strconv.ParseFloat(t.Data, 64)

		if err != nil {
			println("Error converting data to int")
			return 0
		}
		return res
	}
	// if it's a message, for every character in the string, it costs one bitcoin
	return float64(len(t.Data))

}
