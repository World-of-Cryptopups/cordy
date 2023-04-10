package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/World-of-Cryptopups/eosrpc.go"
	"github.com/tbdsux/mini-go/requester"
)

type WhitelistResponseProps struct {
	Wallet        string `json:"wallet"`
	TransactionId string `json:"transaction_id"`
}

var (
	client  = &http.Client{}
	request = requester.NewRequester(client)
)

func AddWhitelist(wallet string) (WhitelistResponseProps, error) {
	var resp WhitelistResponseProps

	body := map[string]interface{}{
		"wallet": wallet,
	}
	encBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/whitelist", WHITELIST_API), bytes.NewBuffer(encBody))
	if err != nil {
		return resp, err
	}

	req.Header.Set("X-Space-App-Key", WHITELIST_API_KEY)

	r, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	err = json.NewDecoder(r.Body).Decode(&resp)

	return resp, err
}

func RemoveWhitelist(wallet string) (WhitelistResponseProps, error) {
	var resp WhitelistResponseProps

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/whitelist/%s", WHITELIST_API, wallet), nil)
	if err != nil {
		return resp, err
	}

	req.Header.Set("X-Space-App-Key", WHITELIST_API_KEY)

	r, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&resp)

	return resp, err
}

func GetWhitelists() ([]string, error) {
	table, err := chain.GetTableRows(eosrpc.GetTableRowsProps{
		Code:  "wocgalleryrw",
		Table: "whiteconfig",
		Scope: "wocgalleryrw",
		Limit: 1,
		JSON:  true,
	})

	if err != nil {
		SendLog(&LogProps{
			Type:        LogTypeError,
			Title:       "Fetch Whitelists",
			Description: "Failed to fetch whitelists from contract. Please contact admin / developer.",
			Message:     fmt.Sprintf("%s", err),
		})

		return []string{}, err
	}

	config, ok := table.Rows[0].(map[string]interface{})
	if !ok {
		return []string{}, errors.New("invalid type map")
	}

	rawWallets := config["wallets"].([]interface{})
	wallets := make([]string, len(rawWallets))
	for i, v := range rawWallets {
		wallets[i] = v.(string)
	}

	return wallets, nil
}
