package admin

import (
	"fmt"
	"strings"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

var GetVerifiedWalletsCommand = &minidis.SlashCommandProps{
	Name:        "get-verified-wallets",
	Description: "Gets the wallets of the verified users.",
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

		users, err := lib.GetAllUser()
		if err != nil {
			return c.Edit("Failed to get all registered users / members.")
		}

		wallets := []string{}
		for _, v := range users {
			wallets = append(wallets, v.Wallet)
		}

		return c.EditC(minidis.EditProps{
			Content: fmt.Sprintf("Here are the wallets of the verified members / users...\n\nTotal registered members: %d \t", len(users)),
			Attachments: []*discordgo.File{{
				ContentType: "text/plain",
				Name:        "wallets.txt",
				Reader:      strings.NewReader(strings.Join(wallets, "\n")),
			}},
		})
	},
}
