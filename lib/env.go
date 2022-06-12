package lib

import (
	"os"
	"strings"
)

var (
	GUILD = strings.Split(os.Getenv("GUILD"), ",")
	TOKEN = os.Getenv("TOKEN")
)
