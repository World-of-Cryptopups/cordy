package lib

import (
	"net/http"

	"github.com/World-of-Cryptopups/atomicassets-go"
)

var Atom, _ = atomicassets.NewCustom("https://atomic.wax.eosrio.io/atomicassets/v1", &http.Client{})
