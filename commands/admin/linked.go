package admin

import (
	"strings"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/utils"
	"github.com/bwmarrin/discordgo"
)

var ListUnreg = &minidis.SlashCommandProps{
	Name:                     "list-unreg",
	Description:              "List members that have `Linked Pup` role but is not registered to the bot.",
	DefaultMemberPermissions: 0,
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		notregMembers := []*discordgo.Member{}
		mems := GetAllMembers(c.GuildId, c.Session)

		for _, v := range mems {
			if utils.Includes(lib.VERIFIED_ROLE, v.Roles) {
				// check if user is registered, if it is do not add to list
				if _, exists := lib.GetUser(v.User.ID); exists {
					continue
				}

				notregMembers = append(notregMembers, v)
			}
		}

		names := []string{}
		for _, v := range notregMembers {
			names = append(names, v.User.String())
		}

		return c.EditC(minidis.EditProps{
			Content: "Here are the members that have a `Linked Pup` role but are not registered (or was accidentally removed) to me...",
			Attachments: []*discordgo.File{
				{
					ContentType: "text/plain",
					Name:        "linked-not-registered.txt",
					Reader:      strings.NewReader(strings.Join(names, "\n")),
				},
			},
		})

	},
}
