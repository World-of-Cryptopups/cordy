package commands

import (
	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

var linkCommand = &minidis.SlashCommandProps{
	Name:        "link",
	Description: "Link your wallet with your discord id",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "wallet",
			Description: "Your wax wallet address",
			Required:    true,
		},
	},
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(true)

		wallet := c.Options["wallet"].StringValue()
		if wallet == "" {
			return c.Reply("wallet cannot be empty")
		}

		userid := c.Author.ID

		usersBase := lib.UsersBase()
		user := lib.User{
			Key:    userid,
			ID:     userid,
			Wallet: wallet,
		}

		if _, err := usersBase.Put(user); err != nil {
			lib.PrintError(err)

			_, err = c.Followup("There was a problem trying to link your wallet, if the problem persists, please contact an admin.")
			return err
		}

		return c.Edit("Successfully linked your wallet with your User ID. You can now check you DPS stats with `/dps` command and your role will be updated.")
	},
}
