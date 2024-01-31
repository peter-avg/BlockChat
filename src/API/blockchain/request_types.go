package blockchain

type RegisterNodeRequest struct {
    IP      string `json:"ip"`
    Port    string `json:"port"`
	Modulus  string `json:"modulus"` 
	Exponent int    `json:"exponent"`
}
