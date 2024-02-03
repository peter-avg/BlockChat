package blockchain

import (
    // "fmt"
    "net/http"
    "crypto/rand"
    "log"
    "bytes"
	"crypto/rsa"
    "encoding/json"
)

// NodeInfo struct contains communication info about other nodes
type NodeInfo struct {
    Id int `json:"id"`
    IP string `json:"IP"`
    PORT string `json:"PORT"`
	PublicKey *rsa.PublicKey `json:"PublicKey"`
    Balance int `json:"Balance"`
}

// Node struct contains Blockchain info
type Node struct {
    Id int `json:"id"`
    Nonce int `json:"nonce"`
    Wallet Wallet `json:"wallet"`
    Chain Blockchain `json:"chain"`
    Ring []NodeInfo `json:"ring"`
    CurrentBlock Block `json:"CurrentBlock"`
    Stake int `json:"stake"`
}

// NewNodeInfo creates and returns a new NodeInfo
func NewNodeInfo(id int, ip string, port string,
                 PublicKey *rsa.PublicKey, balance int) *NodeInfo {
    return &NodeInfo {
        Id: id,
        IP: ip, 
        PORT: port,
        PublicKey: PublicKey,
        Balance: balance,
    }
}

// NewNode creates and returns a new Node
func NewNode(id int, chain Blockchain, ring []NodeInfo) *Node {
    return &Node {
        Id: id,
        Nonce: 0,
        Chain: chain,
        Ring: ring,
    }
}

// JSONify serializes NodeInfo into a JSON string
func (ni *NodeInfo) JSONify() (string, error) {
    jsonBytes, err := json.Marshal(ni)
    return string(jsonBytes), err
}

// JSONify serializes Node into a JSON string
func (n *Node) JSONify() (string, error) {
    jsonBytes, err := json.Marshal(n)
    return string(jsonBytes), err
}

// Creating a new block
func (n *Node) CreateNewBlock() {
    if len(n.Chain.Chain) == 0 {
        new_block := NewBlock(0,"1");
        n.CurrentBlock = *new_block;
    } else { 
        new_block := NewBlock(0,"0");
        n.CurrentBlock = *new_block;
    }
}

// Adding Info for a new Node in Ring
func (n *Node) AddNewInfo(info *NodeInfo) {
    n.Ring = append(n.Ring, *info);
    n.Nonce++;
}

// Generating Wallet for Node
func (n *Node) GenerateWallet() {
    n.Wallet = *NewWallet();
}

// Broadcast transaction to all nodes
func (n *Node) BroadcastTransaction(transaction *Transaction) {
    for _, node := range n.Ring {

        if node.Id != n.Id {
            SendTransaction(transaction, node.IP, node.PORT, node.Id);
        }
    }
}

// Send a transaction to a node
// =============================
func SendTransaction(transaction *Transaction, IP string, PORT string, ID int) { 

    send_address := "http://" + IP + ":" + PORT + "/blockchat_api/receive_transaction";

    if ID != transaction.ReceiverAddress {
        bits := 2048
        privateKey, err := rsa.GenerateKey(rand.Reader, bits)
        if err != nil {
            panic("failed to generate private key")
        }
        pK := &privateKey.PublicKey
        transaction.SenderAddress = pK;
    }

    request_body, err := json.Marshal(transaction);
    if err != nil {
        log.Println(err);
        return
    }

    response, err := http.Post(send_address, "application/json", bytes.NewBuffer(request_body));
    if err != nil { 
        log.Println(err);
        return
    }

    log.Println(response);
    // defer response.Body.Close();

    return
}

