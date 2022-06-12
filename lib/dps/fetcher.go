package dps

import (
	"github.com/World-of-Cryptopups/atomicassets-go"
	"github.com/World-of-Cryptopups/cordy/lib"
)

var a = atomicassets.New()

// FetchAllAssets queries all of the assets in the specified schema.
func FetchAllAssets(wallet string, schema string) []atomicassets.AssetsDataProps {
	allData := []atomicassets.AssetsDataProps{}

	page := 1
	for {
		q, err := a.GetAssets(&atomicassets.GetAssetsQuery{
			CollectionName: "cryptopuppie",
			Limit:          1000,
			Owner:          wallet,
			SchemaName:     schema,
			Page:           page,
		})

		if err != nil {
			lib.LogError(err)
		}

		allData = append(allData, q.Data...)

		if len(q.Data) <= 1000 {
			break
		}

		page++
	}

	return allData
}
