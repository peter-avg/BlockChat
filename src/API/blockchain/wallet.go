package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"strconv"
	"fmt"
)

type Wallet struct {
	PublicKey    *rsa.PublicKey
	PrivateKey   *rsa.PrivateKey
	Balance      int
	NodeID       int
	Nonce        int
}

func NewWallet() *Wallet {
	bits := 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic("failed to generate private key")
	}
	publicKey := &privateKey.PublicKey

	return &Wallet{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		Balance:      0,
		NodeID:       0,
		Nonce:        0,
	}
}

// Add a transaction to the wallet
func (w *Wallet) AddTransaction(newTransaction Transaction) int {
	w.Nonce++
    fmt.Println("Nonce updated to " + strconv.Itoa(w.Nonce));
	return w.Nonce
}

// Sign a transaction (sender)
func (w *Wallet) SignTransaction(transaction *Transaction) []byte {
	hashed := sha256.Sum256([]byte(transaction.Data))
	signature, err := rsa.SignPSS(rand.Reader, w.PrivateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
        fmt.Println("Could not sign signature");
		return nil
	}
	return signature
}

// Verify a transaction (receiver)
func (w *Wallet) VerifyTransaction(transaction *Transaction) bool {
	hashed := sha256.Sum256([]byte(transaction.Data))
	err := rsa.VerifyPSS(w.PublicKey, crypto.SHA256, hashed[:], transaction.Signature, nil)
	return err == nil
}

