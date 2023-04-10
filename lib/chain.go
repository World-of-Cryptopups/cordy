package lib

import "github.com/World-of-Cryptopups/eosrpc.go"

var (
	wax   = eosrpc.New("https://wax.greymass.com")
	chain = wax.NewChainAPI()
)
