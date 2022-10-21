package lib

import (
	"net/http"

	"github.com/TheBoringDude/scuffed-go/requester"
)

type WhitelistResponseProps struct {
	Wallet        string `json:"wallet"`
	TransactionId string `json:"transaction_id"`
}

var request = requester.NewRequester(&http.Client{})

func AddWhitelist(wallet string) (WhitelistResponseProps, error) {
	var resp WhitelistResponseProps

	body := map[string]interface{}{
		"wallet": wallet,
	}
	err := request.Post(WHITELIST_API, body, &resp)

	return resp, err
}
