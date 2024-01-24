package main

import (
	"api/blockchain/blockchain"
    "fmt"
)

func main() {

    // Initialising Blockchain
    chain := blockchain.NewBlockchain();
    genesis_block := blockchain.NewBlock(0,"1");
    chain.AddBlock(*genesis_block);


    node := blockchain.NewNode(1,*chain,nil);
    nodeinfo := blockchain.NewNodeInfo(1,"127.0.0.1","9867", node.Wallet.PublicKey,0);
    node.AddNewInfo(*nodeinfo);

    fmt.Println(node.Ring[0]);
}
