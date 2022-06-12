package commands

import (
	"fmt"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/dps"
)

var dpsCommand = &minidis.SlashCommandProps{
	Name:        "dps",
	Description: "Check your pup collection dps stats",
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(true)

		userid := c.Author.ID

		user, exists := lib.GetUser(userid)
		if !exists {
			return c.Edit("You currently do not have a linked wallet to your acc.")
		}

		data := dps.Calculate(user.Wallet)
		fmt.Println(data)

		return c.Edit("here")
	},
}
