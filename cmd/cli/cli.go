package main

import (
	"block-chat/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// JSON structs for Unmarshall

type ViewResponse struct {
	LastBlock string `json:"last_block"`
}

type LastBlockData struct {
	Index        int      `json:"index"`
	PreviousHash string   `json:"previous_hash"`
	Transactions []string `json:"transactions"`
	Validator    string   `json:"validator"`
	Timestamp    string   `json:"timestamp"`
}

type BalanceResponse struct {
	Balance   float32 `json:"balance"`
	SoftStake float32 `json:"soft_stake"`
}

var portNumber int
var stakeAmount float64

// Cli Function declaration
var port = &cli.IntFlag{
	Name:        "port",
	Usage:       "-{port,-port} <port_number> : To choose a specific port (9921 is the default)",
	Destination: &portNumber,
}
var transaction = &cli.BoolFlag{
	Name:  "t",
	Usage: "-{t,-t} <recipient_id> <Message or Number of BlockChat Coins> : To produce a transaction",
}

var stake = &cli.Float64Flag{
	Name:        "stake",
	Usage:       "-{stake,-stake} <amount> : To produce a stake",
	Destination: &stakeAmount,
}

var view = &cli.BoolFlag{
	Name:  "view",
	Usage: "-{view,-view} : To view last block",
}

var balance = &cli.BoolFlag{
	Name:  "balance",
	Usage: "-{balance,-balance} : To show balance",
}

var help = &cli.BoolFlag{
	Name:  "help",
	Usage: "-{help,-help} : Show available commands",
}

//goland:noinspection SpellCheckingInspection
func main() {

	app := &cli.App{

		Name:  "BlockChat CLI",
		Usage: "Used to interact with the BlockChat Application",

		Flags: []cli.Flag{
			port,
			transaction,
			stake,
			view,
			balance,
		},

		Action: func(c *cli.Context) error {
			if c.IsSet("help") || (c.NArg() == 0 && c.NumFlags() == 0) {
				err := cli.ShowAppHelp(c)
				if err != nil {
					return err
				}
				return nil
			}

			var isTransactionSet bool = c.IsSet("t")
			var isStakeSet bool = c.IsSet("stake")
			var isViewSet bool = c.IsSet("view")
			var isBalanceSet bool = c.IsSet("balance")
			var isPortSet bool = c.IsSet("port")

			var apiUrl string = config.API_URL
			if isPortSet {
				log.Println("Using Port Specified : " + strconv.Itoa(portNumber))
				apiUrl += strconv.Itoa(portNumber)
			} else {
				log.Println("Port not specified! Set to default : " + config.DEFAULT_PORT + ".")
				apiUrl += config.DEFAULT_PORT
			}

			apiUrl += "/blockchat_api/"

			// Make Transaction Function Implementation
			// ========================================
			if isTransactionSet {
				log.Println("txn")
				data := make(map[string]interface{})
				recipientIdStringOld := c.Args().Get(0)
				messageOrBCC := c.Args().Get(1)
				recipientIdString := recipientIdStringOld[2:]

				recipientId, recipientIdConvertToIntError := strconv.Atoi(recipientIdString)
				if recipientIdConvertToIntError != nil {
					log.Println("Error : <recipient_id> must be of type int.\n" + recipientIdConvertToIntError.Error())
				}
				//log.Println("firstParam : ", recipientId)
				//log.Println("secondParam : " + messageOrBCC)

				transactionUrl := apiUrl + "send_transaction"
				if recipientIdString == "" {
					fmt.Println("Usage: -{t,-t} <recipient_address> <Message or Number of BlockChat Coins> : To produce a transaction")
					return nil
				}

				_, err := strconv.ParseFloat(messageOrBCC, 32)

				if err != nil {
					var message string = c.Args().Get(1)
					if message == "" {
						fmt.Println("Usage: -{t,-t} <recipient_address> <Message or Number of BlockChat Coins> : To produce a transaction")
						return nil
					}
					//log.Println("It is a message")

					data["recipient_id"] = recipientId
					data["message_or_bitcoin"] = 0
					data["data"] = message
				}

				if err == nil {
					numberOfBlockChatCoins := messageOrBCC
					log.Println("It is BCC : " + messageOrBCC)

					data["recipient_id"] = recipientId
					data["message_or_bitcoin"] = 1
					data["data"] = numberOfBlockChatCoins
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					return nil
				}

				r, err := http.NewRequest("POST", transactionUrl, bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return nil
				}

				r.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(r)
				if err != nil {
					fmt.Println("Error sending request:", err)
					return nil
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					fmt.Println("Your transaction has been submitted")
				} else {
					fmt.Println("Failed to submit transaction: ", resp.StatusCode)
				}

			}

			// Set Stake Function Implementation
			// =================================
			if isStakeSet {
				requestData := make(map[string]interface{})
				stakeUrl := apiUrl + "set_stake"

				requestData["stake"] = stakeAmount

				jsonData, err := json.Marshal(requestData)
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					return nil
				}

				r, err := http.NewRequest("POST", stakeUrl, bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error creating stake request:", err)
					return nil
				}
				r.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(r)
				if err != nil {
					fmt.Println("Error sending request:", err)
					return nil
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					fmt.Println("Your stake transaction has been submitted")
				} else {
					fmt.Println("Failed to submit stake transaction: ", resp.StatusCode)
				}
			}

			// View Last Block Function Implementation
			// =======================================
			if isViewSet {
				viewUrl := apiUrl + "get_last_block"

				resp, err := http.Get(viewUrl)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				log.Println("Response Body : " + string(body))
				if err != nil {
					log.Fatal(err)
				}
			}

			// View Balance Function Implementation
			// ====================================
			if isBalanceSet {
				balanceUrl := apiUrl + "get_balance"

				resp, err := http.Get(balanceUrl)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					log.Fatal(err)
				}
				var apiResponse BalanceResponse
				if err := json.Unmarshal(body, &apiResponse); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Your Balance is : %.3f BlockChat Coins.\nYour Stake is : %.3f", apiResponse.Balance, apiResponse.SoftStake)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}
