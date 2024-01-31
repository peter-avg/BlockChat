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

