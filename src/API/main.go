package main

import (
    "bytes"
    "encoding/json"
	"api/blockchain/blockchain"
    "io"
	"fmt"
    "log"
	// "fmt"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
)

var MyNode blockchain.Node = blockchain.MyNode;
var BOOTSTRAP_IP string = blockchain.BOOTSTRAP_IP;
var BOOTSTRAP_PORT string = blockchain.BOOTSTRAP_PORT;
var CAPACITY int = blockchain.CAPACITY;

func main() {

    router := gin.Default();

    router.POST("/register_node", blockchain.RegisterNode);

    IP,err := blockchain.GetIP();

    if err != nil {
        log.Fatal("Could not get IP");
    }

    var PORT string;
    var nodes int;
    var bootstrap bool;

    flag.StringVar(&PORT,"p", "5000", "Port to run on");
    flag.IntVar(&nodes,"n", 1, "Number of nodes in chain");
    flag.BoolVar(&bootstrap,"b", false, "If node is bootstrap node");

    flag.Parse();

    if bootstrap {

        // Setup the Bootstrap node
        MyNode.Id = 0;
        MyNode.GenerateWallet();
        log.Println("Bootstrap Node public key: ",MyNode.Wallet.PublicKey);
        MyNodeInfo := blockchain.NewNodeInfo(MyNode.Id, BOOTSTRAP_IP, BOOTSTRAP_PORT, MyNode.Wallet.PublicKey, nodes*1000);
        MyNode.AddNewInfo(MyNodeInfo);

        // Setup the Genesis Block
        GenesisBlock := blockchain.NewBlock(0, "1");
        FirstTransaction := blockchain.NewTransaction(0,0, true, fmt.Sprint(1000*nodes), 1);
        GenesisBlock.AddTransaction(FirstTransaction, CAPACITY);
        GenesisBlock.Hashify();

        // Insert the Genesis Block into the Blockchain
        MyNode.Chain.AddBlock(*GenesisBlock);

        router.Run(BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT);

    } else {

        entry_address := "http://" + BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT + "/register_node";
        MyNode.GenerateWallet();

        requestBody,_ := json.Marshal(map[string]interface{}{
            "ip": IP,
            "port": PORT,
            "modulus": MyNode.Wallet.PublicKey.N,
            "exponent": MyNode.Wallet.PublicKey.E,
        })

        response, err := http.Post(entry_address, "application/json", bytes.NewBuffer(requestBody));

        if err != nil {
            log.Fatal("Could not connect to Bootstrap Node");
        }

        body, err := io.ReadAll(response.Body)
        if err != nil {
           log.Fatalln(err)
        }

        sb := string(body)
        log.Printf(sb)

        router.Run(IP + ":" + PORT);
        
    }

}







