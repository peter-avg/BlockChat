package main

import (
    "bytes"
    "encoding/json"
	"api/blockchain/blockchain"
	"fmt"
    "log"
	// "fmt"
	"flag"
	// "github.com/gin-gonic/gin"
    // "github.com/gin-contrib/cors"
	"net/http"
)

var MyNode blockchain.Node;
var BOOTSTRAP_IP string = blockchain.BOOTSTRAP_IP;
var BOOTSTRAP_PORT string = blockchain.BOOTSTRAP_PORT;
var CAPACITY int = blockchain.CAPACITY;

func main() {

    router := InitRouter();

    // IP,err := blockchain.GetIP();
    // if err != nil {
    //     log.Fatal("Could not get IP");
    // }

    var PORT string;
    var nodes int;
    var bootstrap bool;

    fmt.Println("hello i am a node");

    flag.StringVar(&PORT,"p", "5000", "Port to run on");
    flag.IntVar(&nodes,"n", 3, "Number of nodes in chain");
    flag.BoolVar(&bootstrap,"b", false, "If node is bootstrap node");

    fmt.Println("hello i am a node");

    flag.Parse();

    if bootstrap {

        // Setup the Bootstrap node
        MyNode.Id = 0;
        MyNode.GenerateWallet();
        MyNodeInfo := blockchain.NewNodeInfo(MyNode.Id, BOOTSTRAP_IP, BOOTSTRAP_PORT, MyNode.Wallet.PublicKey, nodes*1000);
        MyNode.Wallet.Balance = nodes*1000;
        MyNode.AddNewInfo(MyNodeInfo);
        log.Println(MyNode.Ring)

        // Setup the Genesis Block
        GenesisBlock := blockchain.Block{
            Index: 0,
            Timestamp: blockchain.GetTimestamp(),
            Transactions: []blockchain.Transaction{},
            Validator: 0,
            CurrentHash: "",
            PreviousHash: "1",
        };

        var FirstTransaction = blockchain.Transaction{
            SenderAddress: MyNode.Wallet.PublicKey,
            ReceiverAddress: 0,
            TypeOfTransaction: true,
            Data : fmt.Sprint(1000*nodes),
            Nonce : 0,
            TransactionID : "",
            Signature : nil,
        };
        GenesisBlock.AddTransaction(FirstTransaction, CAPACITY);
        GenesisBlock.Hashify();

        // Insert the Genesis Block into the Blockchain
        MyNode.Chain.AddBlock(GenesisBlock);

        router.Run(BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT);

    } else {

        entry_address := "http://" + BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT + "/blockchat_api/register_node";
        MyNode.GenerateWallet();

        requestBody,_ := json.Marshal(map[string]interface{}{
            "ip": BOOTSTRAP_IP,
            "port": PORT,
            "modulus": MyNode.Wallet.PublicKey.N,
            "exponent": MyNode.Wallet.PublicKey.E,
        })

        response, err := http.Post(entry_address, "application/json", bytes.NewBuffer(requestBody));
        if err != nil {
            log.Fatal(err);
        }

        MyNode.Chain, MyNode.Ring, MyNode.Wallet.Balance, err = blockchain.DeserializeRegisterNodeResponse(response);
        if err != nil {
            log.Fatal(err);
        }

        fmt.Println(MyNode.Chain);
        fmt.Println(MyNode.Ring);

        router.Run(BOOTSTRAP_IP + ":" + PORT);
        
    }
}
