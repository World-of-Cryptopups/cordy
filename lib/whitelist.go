package lib

import (
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
	err := request.Post(fmt.Sprintf("%s/whitelist", WHITELIST_API), body, &resp)

	return resp, err
}

func RemoveWHitelist(wallet string) (WhitelistResponseProps, error) {
	var resp WhitelistResponseProps

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/whitelist/%s", WHITELIST_API, wallet), nil)
	if err != nil {
		return resp, err
	}

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

func RemoveWhitelist(wallet string) error {

	return nil
}
