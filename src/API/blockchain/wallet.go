package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
func (w *Wallet) AddTransaction() int {
	w.Nonce++
	return w.Nonce
}

// Sign a transaction (sender)
func (w *Wallet) SignTransaction(transaction *Transaction) ([]byte,error) {
	hashed := sha256.Sum256([]byte(transaction.Data))
	signature, err := rsa.SignPSS(rand.Reader, w.PrivateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
        fmt.Println("Could not sign signature");
		return nil,err
	}
	return signature,nil
}

// Verify a transaction (receiver)
func (w* Wallet) VerifyTransaction(transaction *Transaction) (bool,error) {
	hashed := sha256.Sum256([]byte(transaction.Data))
	err := rsa.VerifyPSS(transaction.SenderAddress, crypto.SHA256, hashed[:], transaction.Signature, nil)
    return true,err
}

// Deduct money from the wallet
func (w *Wallet) DeductMoney(Amount int) {
    w.Balance -= Amount;
}

// Add money to the wallet
func (w *Wallet) AddMoney(Amount int) {
    w.Balance += Amount;
}
