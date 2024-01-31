package blockchain

import (
    "os"
    "net"
)

func DeserializeTransaction(json string) {
    /* FIX ME */
}

func DeserializeBlock(json string) {
    /* FIX ME */
}

func DeserializeBlockchain(json string) {
    /* FIX ME */
}

func DeserializeWallet(json string) {
    /* FIX ME */
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

