package lib

import (
	"errors"
	"fmt"

	"github.com/World-of-Cryptopups/eosrpc.go"
)

func GetBlacklists() ([]string, error) {
	table, err := chain.GetTableRows(eosrpc.GetTableRowsProps{
		Code:  "wocgalleryrw",
		Table: "config",
		Scope: "wocgalleryrw",
		Limit: 1,
		JSON:  true,
	})

	if err != nil {
		// send log
		SendLog(&LogProps{
			Type:        LogTypeError,
			Title:       "Fetch Blacklists",
			Description: "Failed to fetch blacklists from contract. Please check this error asap.",
			Message:     fmt.Sprintf(`%s`, err),
		})

		return []string{}, err
	}

	config, ok := table.Rows[0].(map[string]interface{})
	if !ok {
		return []string{}, errors.New("invalid type map")
	}

	rawWallets := config["blacklisted_wallets"].([]interface{})
	wallets := make([]string, len(rawWallets))
	for i, v := range rawWallets {
		wallets[i] = v.(string)
	}

	return wallets, nil
}
