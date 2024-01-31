package main

import (
	"api/blockchain/blockchain"
	"fmt"
    "log"
	// "fmt"
    "os"
    "net"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

var BOOTSTRAP_IP string = "127.0.0.1";
var BOOTSTRAP_PORT string = "5000";
var CAPACITY int = 5;
var MyNode blockchain.Node;

func main() {

    router := gin.Default();

    router.GET("/request_entry", RequestEntry);

    IP,err := GetIP();

    log.Println("IP: ", IP);

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
        MyNodeInfo := blockchain.NewNodeInfo(MyNode.Id, BOOTSTRAP_IP, BOOTSTRAP_PORT, MyNode.Wallet.PublicKey, nodes*1000);
        MyNode.AddNewInfo(MyNodeInfo);

        // Setup the Genesis Block
        GenesisBlock := blockchain.NewBlock(0, "1");
        FirstTransaction := blockchain.NewTransaction(0,0, true, fmt.Sprint(1000*nodes), 1);
        GenesisBlock.AddTransaction(FirstTransaction, CAPACITY);
        GenesisBlock.Hashify();

        // Insert the Genesis Block into the Blockchain
        MyNode.Chain.AddBlock(*GenesisBlock);

    } else {

        // entry_address := "http://" + BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT + "/request_entry";
        // MyNode.GenerateWallet();
        // MyNodeInfo := blockchain.NewNodeInfo(0, IP, PORT, MyNode.Wallet.PublicKey,0);




        
    }

    router.Run(IP + ":" + PORT);
}

func RequestEntry(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "id" : fmt.Sprint(MyNode.Ring[len(MyNode.Ring)-1].Id + 1),
        "blockchain" : MyNode.Chain,
    })
};






