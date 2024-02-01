package blockchain

import (
    "math/big"
)

type RegisterNodeRequest struct {
    IP      string `json:"ip"`
    Port    string `json:"port"`
	Modulus  *big.Int `json:"modulus"` 
	Exponent int    `json:"exponent"`
}

type RegisterNodeResponse struct { 
    Id int `json:"id"`
    Blockchain string `json:"blockchain"`
    Ring string `json:"ring"`
    Balance int `json:"balance"`
}

type SetStakeRequest struct {
    Stake int `json:"stake"`
}

type SendTransactionRequest struct {
    Recipient string `json:"recipient_id"`
    Message_or_Bitcoin int `json:"message_or_bitcoin"`
    Data string `json:"data"`
}

