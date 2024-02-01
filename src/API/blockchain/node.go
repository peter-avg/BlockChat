package blockchain

import (
    // "fmt"
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

// Share the Ring with other Nodes
// func (n *Node) ShareRing() {
//     for _, node := range n.Ring {
//         if node.Id == 0 {
//             continue;
//         }
        // node_address := "http://" + node.IP + ":" + node.PORT + "/update_ring";
        // json_string, err := n.JSONify();
        // if err != nil {
        //     continue;
        // }
//     }
// }




