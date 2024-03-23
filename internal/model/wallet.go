package model

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type Wallet struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	Balance    float64
	NodeID     int
	Nonce      int
}

func NewWallet() *Wallet {
	bits := 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic("failed to generate private key")
	}
	publicKey := &privateKey.PublicKey

	return &Wallet{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Balance:    0,
		NodeID:     0,
		Nonce:      0,
	}
}

// AddTransaction Add a transaction to the wallet
func (w *Wallet) AddTransaction() int {
	w.Nonce++
	return w.Nonce
}

// SignTransaction Sign a transaction (sender)
func (w *Wallet) SignTransaction(transaction *Transaction) ([]byte, error) {
	hashed := sha256.Sum256([]byte(transaction.Data))
	signature, err := rsa.SignPSS(rand.Reader, w.PrivateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		fmt.Println("Could not sign signature")
		return nil, err
	}
	return signature, nil
}

// VerifyTransaction Verify transaction signature (receiver)
func (w *Wallet) VerifySignature(Data string, Signature []byte, SenderAddress *rsa.PublicKey) (bool, error) {
	hashed := sha256.Sum256([]byte(Data))
	err := rsa.VerifyPSS(SenderAddress, crypto.SHA256, hashed[:], Signature, nil)
	return true, err
}

// DeductMoney Deduct money from the wallet
func (w *Wallet) DeductMoney(Amount float64) bool {
	if w.Balance-Amount >= 0 {
		w.Balance -= Amount
		return true
	}
	return false
}

// AddMoney Add money to the wallet
func (w *Wallet) AddMoney(Amount float64) {
	w.Balance += Amount
}
