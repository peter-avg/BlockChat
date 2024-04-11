package main

import (
	"block-chat/internal/config"
	"block-chat/internal/endpoints"
	"block-chat/internal/model"
	"block-chat/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	// "fmt"
	"flag"
	// "github.com/gin-gonic/gin"
	// "github.com/gin-contrib/cors"
	"net/http"
)

var MyNode model.Node

var BOOTSTRAP_IP string = config.BOOTSTRAP_IP
var BOOTSTRAP_PORT string = config.BOOTSTRAP_PORT
var CAPACITY int = config.CAPACITY

func main() {
	router := endpoints.InitRouter(&MyNode)

	// IP,err := blockchain.GetIP();
	// if err != nil {
	//     log.Fatal("Could not get IP");
	// }

	var PORT string
	var nodes int
	var bootstrap bool

	flag.StringVar(&PORT, "p", "9921", "Port to run on")
	flag.IntVar(&nodes, "n", 5, "Number of nodes in chain")
	flag.BoolVar(&bootstrap, "b", false, "If node is bootstrap node")

	flag.Parse()

	if bootstrap {

		// Setup the Bootstrap node
		MyNode.Id = 0
		MyNode.GenerateWallet()
		MyNodeInfo := model.NewNodeInfo(MyNode.Id, BOOTSTRAP_IP, BOOTSTRAP_PORT, MyNode.Wallet.PublicKey, 0)
		//MyNode.Wallet.Balance = float64(nodes * 1000)
		MyNode.Nonce = -1
		MyNode.AddNewInfo(MyNodeInfo)
		log.Println(MyNode.Ring)

		// Setup the Genesis Block
		GenesisBlock := model.Block{
			Index:        0,
			Timestamp:    model.GetTimestamp(),
			Transactions: []model.Transaction{},
			Validator:    0,
			CurrentHash:  "",
			PreviousHash: "1",
		}

		var FirstTransaction = model.Transaction{
			SenderAddress:     &config.STAKE_PUBLIC_ADDRESS,
			ReceiverAddress:   MyNode.Wallet.PublicKey,
			TypeOfTransaction: true,
			Data:              fmt.Sprint(1000 * nodes),
			Nonce:             0,
			TransactionID:     "",
			Signature:         nil,
		}
		GenesisBlock.AddTransaction(FirstTransaction, &MyNode)
		GenesisBlock.Hashify()
		GenesisBlock.CurrentHash = "GENESIS_BLOCK"
		// Insert the Genesis Block into the Blockchain
		MyNode.Chain.ValidateBlock(&GenesisBlock, &MyNode)
		log.Println("Last Block : " + MyNode.Chain.GetLastBlock().String())
		router.Run(BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT)

	} else {
		entry_address := "http://" + BOOTSTRAP_IP + ":" + BOOTSTRAP_PORT + "/blockchat_api/register_node"
		MyNode.GenerateWallet()

		requestBody, _ := json.Marshal(map[string]interface{}{
			"ip":       BOOTSTRAP_IP,
			"port":     PORT,
			"modulus":  MyNode.Wallet.PublicKey.N,
			"exponent": MyNode.Wallet.PublicKey.E,
		})

		response, err := http.Post(entry_address, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode == http.StatusBadRequest {
			log.Println(utils.GetErrorMessageFromResponse(response))
			return
		}
		MyNode.Id, MyNode.Chain, MyNode.Ring, MyNode.Wallet.Balance, MyNode.CurrentBlock, err = utils.DeserializeRegisterNodeResponse(response)
		MyNode.Wallet.Nonce = MyNode.CurrentBlock.Transactions[len(MyNode.CurrentBlock.Transactions)-1].Nonce
		if err != nil {
			log.Fatal(err)
		}
		log.Println("MyNode : ", MyNode.String())
		router.Run(BOOTSTRAP_IP + ":" + PORT)
	}
}
