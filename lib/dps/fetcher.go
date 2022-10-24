package dps

import (
	"fmt"
	"log"

	"github.com/World-of-Cryptopups/atomicassets-go"
	"github.com/World-of-Cryptopups/cordy/lib"
)

// FetchAllAssets queries all of the assets in the specified schema.
func FetchAllAssets(wallet string, schema string) []atomicassets.AssetsDataProps {
	allData := []atomicassets.AssetsDataProps{}

	page := 1
	for {
		q, err := lib.Atom.GetAssets(&atomicassets.GetAssetsQuery{
			CollectionName:    "cryptopuppie",
			Limit:             1000,
			Owner:             wallet,
			SchemaName:        schema,
			Page:              page,
			TemplateBlacklist: []string{"613110", "612990"}, // ignore the woof coins
		})

		if err != nil {
			log.Println(err)

			// send log
			lib.SendLog(&lib.LogProps{
				Type:        lib.LogTypeError,
				Title:       "Error Fetching Assets",
				Description: fmt.Sprintf("Failed in fetching the user's assets to calculate the dps, wallet: %s", wallet),
				Message:     fmt.Sprintf("`%v`", err),
			})

			// try to redo
			return FetchAllAssets(wallet, schema)
		}

		allData = append(allData, q.Data...)

		if len(q.Data) <= 1000 {
			break
		}

		page++
	}

	return allData
}
