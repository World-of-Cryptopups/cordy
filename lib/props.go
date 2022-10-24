package lib

type User struct {
	Key       string `json:"key"`
	ID        string `json:"id"`
	Wallet    string `json:"wallet"`
	Token     string `json:"token"`
	IsStopped bool   `json:"is_stopped"`
}

type UserDpsProps struct {
	Key       string   `json:"key"`
	ID        string   `json:"id"`
	Dps       DPSProps `json:"dps"`
	Wallet    string   `json:"wallet"`
	IsStopped bool     `json:"is_stopped"`
}

type WebLoginUserProps struct {
	Key    string `json:"key"`
	Token  string `json:"token"`
	Type   string `json:"type"`
	Wallet string `json:"wallet"`
	Linked bool   `json:"linked,omitempty"`
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
