package main

import (
	"api/blockchain/blockchain"
	"fmt"
)

func main() {
    var wallet = blockchain.NewWallet();

    var block = blockchain.NewBlock(1,nil,1,"");
    
    var transaction = blockchain.NewTransaction(0,0,true,"",0,0);

    fmt.Println(block);
    fmt.Println(transaction);

    transaction.Nonce = wallet.AddTransaction(*transaction);
    transaction.Signature = wallet.SignTransaction(transaction);

    if (wallet.VerifyTransaction(transaction)) {
        fmt.Println(transaction);
    }
}
