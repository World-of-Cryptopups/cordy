package lib

type User struct {
	Key    string `json:"key"`
	ID     string `json:"id"`
	Wallet string `json:"wallet"`
}

type UserDpsProps struct {
	Key string   `json:"key"`
	ID  string   `json:"id"`
	Dps DPSProps `json:"dps"`
}

type DPSProps struct {
	PupSkinCards int           `json:"pupskincards"`
	PuppyCards   int           `json:"puppycards"`
	PupItems     DPSItemsProps `json:"pupitems"`
}

type DPSItemsProps struct {
	Real int `json:"real"`
	Raw  int `json:"raw"`
}
