package admin

import (
	"strings"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/tbdsux/mini-go/mini"
)

var ListNotExistent = &minidis.SlashCommandProps{
	Name:                     "list-nonexistent",
	Description:              "List members that left the server but still exists in the whitelist.",
	DefaultMemberPermissions: 0,
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		whitelists, err := lib.GetWhitelists()
		if err != nil {
			return err
		}

		users, err := lib.GetAllUser()
		if err != nil {
			return err
		}

		usersWallets := []string{}
		for _, v := range users {
			usersWallets = append(usersWallets, v.Wallet)
		}

		// not exists in server but is whitelisted
		notExistentWallets := []string{}

		for _, v := range whitelists {
			if !mini.Exists(usersWallets, v) {
				notExistentWallets = append(notExistentWallets, v)
			}
		}

		return c.EditC(minidis.EditProps{
			Content: "Wallets that are whitelisted but have left the Discord Server.",
			Attachments: []*discordgo.File{
				{
					ContentType: "text/plain",
					Name:        "nonexistent-whitelist.text",
					Reader:      strings.NewReader(strings.Join(notExistentWallets, "\n")),
				},
			},
		})
	},
}

// INFO: currently too much work
//
// var RemoveNonExistent = &minidis.SlashCommandProps{
// 	Name:                     "remove-nonexistent",
// 	Description:              "Remove members that left the server but still exists in the whitelist.",
// 	DefaultMemberPermissions: 0,
// 	Execute: func(c *minidis.SlashContext) error {
// 		c.DeferReply(false)

// 		return nil
// 	},
// }
