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
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		// https://github.com/bwmarrin/discordgo/issues/1024
		perms, err := c.Session.UserChannelPermissions(c.Author.ID, c.ChannelId)
		if err != nil {
			_, err := c.Followup("There was a problem getting the user's permissions.")
			return err
		}
		if perms&discordgo.PermissionAdministrator == 0 {
			// not admin
			return c.Edit("You do not have permission to perform such actions!")
		}

		mentioned := c.Options["user"].UserValue(c.Session)

		user, exists := lib.GetUser(mentioned.ID)
		if !exists {
			return c.Edit("The user is currently not registered to me.")
		}

		// remove the user from the databases
		if err := lib.UpdateUser(user.ID, user.Wallet); err != nil {
			return c.Edit(fmt.Sprintf("Error: %v", err))
		}

		return c.Edit(fmt.Sprintf("Successfully reset the linked account of **%s**", mentioned.String()))
	},
}
