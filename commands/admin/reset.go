package admin

import (
	"fmt"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

var ResetAccountCommand = &minidis.SlashCommandProps{
	Name:        "reset",
	Description: "Reset the linked account of a user.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "The user to reset the account",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	},
	DefaultMemberPermissions: 0,
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		mentioned := c.Options["user"].UserValue(c.Session)

		user, exists := lib.GetUser(mentioned.ID)
		if !exists {
			return c.Edit("The user is currently not registered to me.")
		}

		// remove the user from the databases
		if err := lib.UnlinkUser(user.ID, user.Wallet, lib.UserSession{
			UserExists: true,
			Session:    c.Session,
			GuildID:    c.GuildId,
		}, "account has been reset"); err != nil {
			return c.Edit(fmt.Sprintf("Error: %v", err))
		}

		return c.Edit(fmt.Sprintf("Successfully reset the linked account of **%s**", mentioned.String()))
	},
}
