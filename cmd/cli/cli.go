package cli

import (
	"block-chat/internal/config"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	Balance float32 `json:"balance"`
}

// Cli Function declaration

var transaction = &cli.BoolFlag{
	Name:  "t",
	Usage: "-{t,-t} <recipient_address> <message or bitcoin data> : To produce a transaction",
}

var stake = &cli.BoolFlag{
	Name:  "stake",
	Usage: "-{stake,-stake} <amount> : To produce a stake",
}

var view = &cli.BoolFlag{
	Name:  "view",
	Usage: "-{view,-view} : To view last block",
}

var balance = &cli.BoolFlag{
	Name:  "balance",
	Usage: "-{balance,-balance} : To show balance",
}

//goland:noinspection SpellCheckingInspection
func AppCLI() {

	var apiUrl string = config.API_URL

	app := &cli.App{

		Name:  "BlockChat",
		Usage: "Used to interact with the BlockChat Application",

		Flags: []cli.Flag{
			transaction,
			stake,
			view,
			balance,
		},

		Action: func(c *cli.Context) error {

			var transaction bool = c.Bool("t")
			var stake bool = c.Bool("s")
			var view bool = c.Bool("v")
			var balance bool = c.Bool("b")

			if !transaction && !stake && !view && !balance {

			}
			// Make Transaction Function Implementation
			// ========================================
			if transaction {

				recipientId := c.Args().Get(0)
				transactionUrl := apiUrl + "send_transaction"
				data := url.Values{}

				if recipientId == "" {
					fmt.Println("Usage: -{t,-t} <recipient_address> <message or bitcoin data> : To produce a transaction")
					return nil
				}

				_, err := strconv.ParseFloat(c.Args().Get(1), 32)

				if err != nil {

					var message string = c.Args().Get(1)

					if message == "" {
						fmt.Println("Usage: -{t,-t} <recipient_address> <message or bitcoin data> : To produce a transaction")
						return nil
					}

					data.Set("recipient_id", recipientId)
					data.Set("message_or_bitcoin", "0")
					data.Set("data", message)

				}

				if err == nil {
					bc := c.Args().Get(1)

					data.Set("recipient_id", recipientId)
					data.Set("message_or_bitcoin", "1")
					data.Set("data", bc)
				}

				r, err := http.NewRequest("POST", transactionUrl, strings.NewReader(data.Encode()))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return nil
				}

				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
			if stake {

				stakeUrl := apiUrl + "set_stake"

				_, err := strconv.ParseFloat(c.Args().Get(0), 32)

				if err != nil {
					fmt.Println("Usage: cli -stake <amount> : To produce a stake")
					return nil
				}

				stakeValue := c.Args().Get(0)

				data := url.Values{}
				data.Set("stake", stakeValue)

				r, err := http.NewRequest("POST", stakeUrl, strings.NewReader(data.Encode()))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return nil
				}

				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

				client := &http.Client{}
				resp, err := client.Do(r)
				if err != nil {
					fmt.Println("Error sending request:", err)
					return nil
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					fmt.Println("Stake was set")
				} else {
					fmt.Println("Failed to set stake, status code:", resp.StatusCode)
				}
			}

			// View Last Block Function Implementation
			// =======================================
			if view {
				viewUrl := apiUrl + "get_last_block"

				resp, err := http.Get(viewUrl)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					log.Fatal(err)
				}

				var apiResponse ViewResponse
				if err := json.Unmarshal(body, &apiResponse); err != nil {
					log.Fatal(err)
				}

				var lastBlock LastBlockData
				if err := json.Unmarshal([]byte(apiResponse.LastBlock), &lastBlock); err != nil {
					log.Fatal(err)
				}

				prettyLastBlock, err := json.MarshalIndent(lastBlock, "", "  ")
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("The Last Block is:")
				fmt.Println(string(prettyLastBlock))
			}

			// View Balance Function Implementation
			// ====================================
			if balance {

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

				fmt.Println("Your Balance is:")
				fmt.Println(apiResponse.Balance)
			}

			return nil
		},
	}

	app.Run(os.Args)
}
