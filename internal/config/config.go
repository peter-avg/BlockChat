package config

import (
	"crypto/rsa"
	"math/big"
)

var BOOTSTRAP_IP string = "127.0.0.1"

var BOOTSTRAP_PORT string = "5000"

var CAPACITY int = 5

var API_URL string = "http://127.0.0.1:"

var DEFAULT_PORT string = "5000"

var STAKE_PUBLIC_ADDRESS rsa.PublicKey = rsa.PublicKey{
	N: big.NewInt(0),
	E: 0,
}
