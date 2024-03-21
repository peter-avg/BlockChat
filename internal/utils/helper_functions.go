package utils

import (
	"block-chat/internal/model"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
)

// Receive a RegisterNodeResponse and convert to Ring and Chain
// ============================================================
func DeserializeRegisterNodeResponse(response *http.Response) (model.Blockchain, []model.NodeInfo, int, error) {
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return model.Blockchain{}, []model.NodeInfo{}, 0, err
	}

	var response_data model.RegisterNodeResponse
	if err := json.Unmarshal(body, &response_data); err != nil {
		return model.Blockchain{}, []model.NodeInfo{}, 0, err
	}

	var blockchain model.Blockchain

	if err := json.Unmarshal([]byte(response_data.Blockchain), &blockchain); err != nil {
		return model.Blockchain{}, []model.NodeInfo{}, 0, err
	}

	var ring []model.NodeInfo
	if err := json.Unmarshal([]byte(response_data.Ring), &ring); err != nil {
		return model.Blockchain{}, []model.NodeInfo{}, 0, err
	}

	return blockchain, ring, response_data.Balance, nil

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
