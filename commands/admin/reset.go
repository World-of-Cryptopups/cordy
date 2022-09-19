package admin

import (
	"fmt"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/deta/deta-go/service/base"
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

		usersBase := lib.UsersBase()
		dpsBase := lib.UsersDpsBase()
		loginsBase := lib.WebLoginBase()

		// remove user from usersbase
		if err = usersBase.Delete(user.ID); err != nil {
			return c.Edit("Failed to remove user data from database." + fmt.Sprintf(" [user:%s]", user.ID))
		}

		// remove user's dps from dpsbase
		if err = dpsBase.Delete(user.ID); err != nil {
			return c.Edit("Failed to remove user's dps data from database." + fmt.Sprintf(" [user:%s]", user.ID))
		}

		// update linked bool in weblogin token basis
		if err = loginsBase.Update(user.Wallet, base.Updates{
			"linked": false,
		}); err != nil {
			return c.Edit("Failed to update user's web login data." + fmt.Sprintf(" [user:%s]", user.ID))
		}

		return c.Edit(fmt.Sprintf("Successfully reset the linked account of **%s**", mentioned.String()))
	},
}
