package lib

import (
	"os"
	"strings"
)

var (
	GUILD = strings.Split(os.Getenv("GUILD"), ",")
	TOKEN = os.Getenv("TOKEN")
	DEV   = os.Getenv("DEV") == "true"

	ADVENTURE_ROLE = os.Getenv("ADVENTURE_ROLE")
	InitRoles      = strings.Split(os.Getenv("ROLES"), ",")
	ALL_ROLES      = append(InitRoles, ADVENTURE_ROLE)
	MOD_ROLE       = os.Getenv("MOD_ROLE")
	ADMIN_ROLE     = os.Getenv("ADMIN_ROLE")
	VERIFIED_ROLE  = os.Getenv("VERIFIED_ROLE")
)
