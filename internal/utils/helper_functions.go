package utils

import (
	"block-chat/internal/model"
	"crypto/rsa"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

// Receive a RegisterNodeResponse and convert to Ring and Chain
// ============================================================
func DeserializeRegisterNodeResponse(response *http.Response) (int, model.Blockchain, []model.NodeInfo, float64, model.Block, error) {
	body, err := io.ReadAll(response.Body)
	log.Println(string(body))
	if err != nil {
		return -1, model.Blockchain{}, []model.NodeInfo{}, 0, model.Block{}, err
	}

	var response_data model.RegisterNodeResponse
	if err := json.Unmarshal(body, &response_data); err != nil {
		log.Println("Line 27")

		return -1, model.Blockchain{}, []model.NodeInfo{}, 0, model.Block{}, err
	}

	var blockchain model.Blockchain
	if err := json.Unmarshal([]byte(response_data.Blockchain), &blockchain); err != nil {
		log.Println("Line 32")

		return response_data.Id, model.Blockchain{}, []model.NodeInfo{}, 0, model.Block{}, err
	}

	var ring []model.NodeInfo
	if err := json.Unmarshal([]byte(response_data.Ring), &ring); err != nil {
		log.Println("Line 37")

		return response_data.Id, blockchain, []model.NodeInfo{}, 0, model.Block{}, err
	}

	var currentBlock model.Block
	if err := json.Unmarshal([]byte(response_data.CurrentBlock), &currentBlock); err != nil {
		log.Println("Line 42")
		return response_data.Id, model.Blockchain{}, ring, 0, model.Block{}, err
	}

	log.Println("response_data.id : " + strconv.Itoa(response_data.Id))
	return response_data.Id, blockchain, ring, response_data.Balance, currentBlock, nil
}

func GetErrorMessageFromResponse(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "Could not read response.Body"
	}

	var responseMap map[string]interface{}

	if err := json.Unmarshal(body, &responseMap); err != nil {
		return "Could not unmarshal error"
	}

	// Access the "error" field
	if errorMsg, ok := responseMap["error"].(string); ok {
		return "Error : " + errorMsg
	} else {
		return "Could not find error"
	}
}

// Get the IP address of the current node
// ======================================
func GetIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return " ", err
	}

	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return " ", err
	}

	for _, addr := range addresses {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return " ", err
}

func FindPublicAddress(nodeInfoArray []model.NodeInfo, nodeId int) *rsa.PublicKey {
	for _, nodeInfo := range nodeInfoArray {
		if nodeInfo.Id == nodeId {
			return nodeInfo.PublicKey
		}
	}
	return nil
}

func Float64ToString(float float64) string {
	return strconv.FormatFloat(float, 'f', -1, 64)
}
