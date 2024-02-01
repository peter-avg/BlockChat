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

