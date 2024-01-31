package blockchain

import (
    "os"
    "net"
    "encoding/json"
)

// JSON to Transaction
func DeserializeTransaction(jsonData string) (Transaction, error) {

    var t Transaction;
    err := json.Unmarshal([]byte(jsonData), &t)
    if err != nil {
        return t, err
    }

    return t, nil
}

// JSON to Block
func DeserializeBlock(jsonData string) (Block, error) {
    /* FIX ME */
    var b Block;

    err := json.Unmarshal([]byte(jsonData), &b)
    if err != nil {
        return b, err
    }

    return b, nil
}

// JSON to Blockchain
func DeserializeBlockchain(jsonData string) (Blockchain, error) {
    /* FIX ME */
    var bc Blockchain;

    err := json.Unmarshal([]byte(jsonData), &bc)
    if err != nil {
        return bc, err
    }

    return bc, nil
}

// JSON to Wallet
func DeserializeWallet(jsonData string) (Wallet, error) {
    /* FIX ME */
    var w Wallet;

    err := json.Unmarshal([]byte(jsonData), &w)
    if err != nil {
        return w, err
    }

    return w, nil
}

// Get the IP address of the current node
func GetIP() (string,error) {
	hostname,err := os.Hostname()
    if err != nil {
        return " ", err
    }

	addresses,err := net.LookupIP(hostname)
    if err != nil {
        return " ",err
    }

	for _, addr := range addresses {
		if ipv4 := addr.To4(); ipv4 != nil {
            return ipv4.String(),nil;
		}
	}

    return " ",err
}

