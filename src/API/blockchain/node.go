package blockchain

import (
    // "fmt"
    "errors"
    "strconv"
    "net/http"
    "sync"
    // "crypto/rand"
    "log"
    "bytes"
    "crypto/rsa"
    "encoding/json"
)

// NodeInfo struct contains communication info about other nodes
// ============================================================
type NodeInfo struct {
    Id int `json:"id"`
    IP string `json:"IP"`
    PORT string `json:"PORT"`
    PublicKey *rsa.PublicKey `json:"PublicKey"`
    Balance int `json:"Balance"`
    Stake int `json:"stake"`
}

// Node struct contains Blockchain info
// ===================================
type Node struct {
    Id int `json:"id"`
    Nonce int `json:"nonce"`
    Wallet Wallet `json:"wallet"`
    Chain Blockchain `json:"chain"`
    Ring []NodeInfo `json:"ring"`
    CurrentBlock Block `json:"CurrentBlock"`
}

// NewNodeInfo creates and returns a new NodeInfo
// =============================================
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
// =====================================
func NewNode(id int, chain Blockchain, ring []NodeInfo) *Node {
    return &Node {
        Id: id,
        Nonce: 0,
        Chain: chain,
        Ring: ring,
    }
}

// JSONify serializes NodeInfo into a JSON string
// =============================================
func (ni *NodeInfo) JSONify() (string, error) {
    jsonBytes, err := json.Marshal(ni)
    return string(jsonBytes), err
}

// JSONify serializes Node into a JSON string
// =========================================
func (n *Node) JSONify() (string, error) {
    jsonBytes, err := json.Marshal(n)
    return string(jsonBytes), err
}

// Creating a new block
// ====================
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
// ==================================
func (n *Node) AddNewInfo(info *NodeInfo) {
    n.Ring = append(n.Ring, *info);
    n.Nonce++;
}

// Generating Wallet for Node
// ==========================
func (n *Node) GenerateWallet() {
    n.Wallet = *NewWallet();
}

// Broadcast new node to all nodes
// ================================
func (n *Node) BroadcastNewNode(info *NodeInfo) {
    for _, node := range n.Ring {
        if node.Id != n.Id && node.Id != info.Id {
            if n.SendNewNode(info, node.IP, node.PORT, node.Id) == false {
                log.Println("Error sending new node to Node ", node.Id);
            }
        }
    }
}

// Send a new node to a node
// ==========================
func (n *Node) SendNewNode(info *NodeInfo, IP string, PORT string, ID int) bool {
    send_address := "http://" + IP + ":" + PORT + "/blockchat_api/receive_new_node";

    request_body, err := json.Marshal(info);
    if err != nil {
        log.Println(err);
        return false
    }

    response, err := http.Post(send_address, "application/json", bytes.NewBuffer(request_body));
    if err != nil { 
        log.Println(err);
        return false
    }

    if response.StatusCode == 200 {
        log.Println("New node sent to Node ", ID);
        return true
    }

    log.Println("New node failed to send to Node ", ID);
    return false
}


// Broadcast transaction to all nodes
// =================================
func (n *Node) BroadcastTransaction(transaction *Transaction) bool {
    var wg sync.WaitGroup
    errChV := make(chan error, len(n.Ring))
    errChS := make(chan error, len(n.Ring))

    for _, node := range n.Ring {
        if node.Id != n.Id {
            wg.Add(1)
            
            go func(node NodeInfo) {
                defer wg.Done()
                if !n.ValidateTransaction(transaction, node.IP, node.PORT, node.Id) {
                    errChV <- errors.New("Validation failed for Node " + strconv.Itoa(node.Id))
                }
            }(node)
        }
    }

    go func() {
        wg.Wait()
        close(errChV)
    }()

    for err := range errChV {
        log.Println(err)
        return false
    }

    for _, node := range n.Ring {
        if node.Id != n.Id {
            wg.Add(1)
            
            go func(node NodeInfo) {
                defer wg.Done()
                if !n.SendTransaction(transaction, node.IP, node.PORT, node.Id) {
                    errChS <- errors.New("Sending failed for Node " + strconv.Itoa(node.Id))
                }
            }(node)
        }
    }

    go func() {
        wg.Wait()
        close(errChS)
    }()

    for err := range errChS {
        log.Println(err)
        return false
    }

    return true
}

// Validate a transaction
// ======================
func (n *Node) ValidateTransaction(transaction *Transaction, IP string, PORT string, ID int) bool {

    send_address := "http://" + IP + ":" + PORT + "/blockchat_api/validate_transaction";

    request_body, err := json.Marshal(transaction);
    if err != nil {
        log.Println(err);
        return false
    }

    response, err := http.Post(send_address, "application/json", bytes.NewBuffer(request_body));
    if err != nil {
        log.Println(err);
        return false;
    }

    if response.StatusCode == 200 {
        log.Println("Transaction validated by Node ", ID);
        return true;
    }

    return false;
}

// Send a transaction to a node
// =============================
func (n *Node) SendTransaction(transaction *Transaction, IP string, PORT string, ID int) bool { 

    send_address := "http://" + IP + ":" + PORT + "/blockchat_api/receive_transaction";

    request_body, err := json.Marshal(transaction);
    if err != nil {
        log.Println(err);
        return false
    }

    response, err := http.Post(send_address, "application/json", bytes.NewBuffer(request_body));
    if err != nil { 
        log.Println(err);
        return false
    }

    if response.StatusCode == 200 {
        log.Println("Transaction sent to Node ", ID);
        return true
    }

    log.Println("Transaction failed to send to Node ", ID);
    return false
}

// Broadcast stake to all nodes
// ============================
func (n *Node) BroadcastStake(info NodeInfo) bool {
    var wg sync.WaitGroup
    errChV := make(chan error, len(n.Ring))
    errChS := make(chan error, len(n.Ring))

    for _, node := range n.Ring {
        if node.Id != n.Id {
            wg.Add(1)
            
            go func(node NodeInfo) {
                defer wg.Done()
                if !n.ValidateStake(node, node.IP, node.PORT, node.Id) {
                    errChV <- errors.New("Validation failed for Node " + strconv.Itoa(node.Id))
                }
            }(node)
        }
    }

    go func() {
        wg.Wait()
        close(errChV)
    }()

    for err := range errChV {
        log.Println(err)
        return false
    }

    for _, node := range n.Ring {
        if node.Id != n.Id {
            wg.Add(1)
            
            go func(node NodeInfo) {
                defer wg.Done()
                if !n.SendStake(node, node.IP, node.PORT, node.Id) {
                    errChS <- errors.New("Sending failed for Node " + strconv.Itoa(node.Id))
                }
            }(node)
        }
    }

    go func() {
        wg.Wait()
        close(errChS)
    }()

    for err := range errChS {
        log.Println(err)
        return false
    }

    return true
}

// Validate a stake
// ================
func (n *Node) ValidateStake(info NodeInfo, IP string, PORT string, ID int) bool {

    send_address := "http://" + IP + ":" + PORT + "/blockchat_api/validate_stake";

    request_body, err := json.Marshal(info);
    if err != nil {
        log.Println(err);
        return false
    }

    response, err := http.Post(send_address, "application/json", bytes.NewBuffer(request_body));
    if err != nil {
        log.Println(err);
        return false;
    }

    if response.StatusCode == 200 {
        log.Println("Stake validated by Node ", ID);
        return true;
    }

    return false;
}

// Send a stake to a node
// ======================
func (n *Node) SendStake(info NodeInfo, IP string, PORT string, ID int) bool {

    send_address := "http://" + IP + ":" + PORT + "/blockchat_api/receive_stake";

    request_body, err := json.Marshal(info);
    if err != nil {
        log.Println(err);
        return false
    }

    response, err := http.Post(send_address, "application/json", bytes.NewBuffer(request_body));
    if err != nil { 
        log.Println(err);
        return false
    }

    if response.StatusCode == 200 {
        log.Println("Stake sent to Node ", ID);
        return true
    }

    log.Println("Stake failed to send to Node ", ID);
    return false
}



