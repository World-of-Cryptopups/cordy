package dps

import (
	"github.com/World-of-Cryptopups/atomicassets-go"
	"github.com/World-of-Cryptopups/cordy/lib"
)

var a = atomicassets.New()

func FetchAllAssets(wallet string, schema string) []atomicassets.AssetsDataProps {
	allData := []atomicassets.AssetsDataProps{}

	for {
		q, err := a.GetAssets(&atomicassets.GetAssetsQuery{
			CollectionName: "cryptopuppie",
			Limit:          1000,
			Owner:          wallet,
			SchemaName:     schema,
		})

		if err != nil {
			lib.LogError(err)
		}

		if len(q.Data) == 0 {
			break
		}

		allData = append(allData, q.Data...)
	}

	return allData
}
