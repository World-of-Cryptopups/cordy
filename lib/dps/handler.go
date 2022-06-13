package dps

import (
	"strconv"
	"strings"

	"github.com/World-of-Cryptopups/atomicassets-go"
	"github.com/World-of-Cryptopups/cordy/lib"
)

type GetDPSProps struct {
}

func cardsDps(wallet string) int {
	data := FetchAllAssets(wallet, "puppycards")
	dps := 0

	for _, v := range data {
		num, _ := strconv.Atoi(v.Data["DPS"].(string))
		dps += num
	}

	return dps
}

func skinsDps(wallet string) ([]atomicassets.AssetsDataProps, int) {
	data := FetchAllAssets(wallet, "pupskincards")
	dps := 0

	for _, v := range data {
		num, _ := strconv.Atoi(v.Data["DPS"].(string))
		dps += num
	}

	return data, dps
}

var (
	_demon = []string{"Demon Queen", "Demon Ace", "Demon King"}
	_mecha = []string{"Mecha Glitter", "Mecha Apollo", "Mecha Draco"}
)

func includes(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

func itemsDps(wallet string, skins []atomicassets.AssetsDataProps) (int, int) {
	data := FetchAllAssets(wallet, "pupitems")

	realDps := 0
	rawDps := 0

	for _, v := range data {
		for _, x := range skins {
			name := strings.TrimSpace(x.Data["name"].(string))

			if includes(_demon, name) {
				name = "Demon"
			}
			if includes(_mecha, name) {
				name = "Mecha"
			}

			itemOwner := v.Data["Item Owner"].(string)
			if name == itemOwner {
				num, _ := strconv.Atoi(v.Data["DPS"].(string))
				realDps += num
				break
			}
		}
	}

	for _, v := range data {
		num, _ := strconv.Atoi(v.Data["DPS"].(string))
		rawDps += num
	}

	return realDps, rawDps
}

// fetches and calculates each type's dps
func Calculate(wallet string) lib.DPSProps {
	cards := cardsDps(wallet)
	skinsData, skins := skinsDps(wallet)
	itemsReal, itemsRaw := itemsDps(wallet, skinsData)

	return lib.DPSProps{
		PupSkinCards: skins,
		PuppyCards:   cards,
		PupItems: lib.DPSItemsProps{
			Real: itemsReal,
			Raw:  itemsRaw,
		},
	}
}
