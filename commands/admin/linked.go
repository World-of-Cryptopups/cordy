package admin

import (
	"fmt"
	"strings"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/tbdsux/mini-go/mini"
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
			if mini.Exists(v.Roles, lib.VERIFIED_ROLE) {
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

var UnlinkUnreg = &minidis.SlashCommandProps{
	Name:                     "unlink-unreg",
	Description:              "Remove members that have a `Linked Pup` role that are not registered to the bot.",
	DefaultMemberPermissions: 0,
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		notregMembers := []*discordgo.Member{}
		mems := GetAllMembers(c.GuildId, c.Session)

		for _, v := range mems {
			if mini.Exists(v.Roles, lib.VERIFIED_ROLE) {
				// check if user is registered, if it is do not add to list
				if _, exists := lib.GetUser(v.User.ID); exists {
					continue
				}

				notregMembers = append(notregMembers, v)
			}
		}

		names := []string{}
		for _, v := range notregMembers {
			if err := c.Session.GuildMemberRoleRemove(c.GuildId, v.User.ID, lib.VERIFIED_ROLE); err != nil {
				fmt.Printf("Failed to remove `Linked Pup` role from user: **`%s`**\n", v.User.String())

				// send log
				lib.SendLog(&lib.LogProps{
					Type:        lib.LogTypeError,
					Title:       "Unlink Error",
					Description: fmt.Sprintf("Failed to remove `Linked Pup` role from user: **`%s`**", v.User.String()),
					Message:     fmt.Sprintf("`%v`", err),
				})
			}

			names = append(names, v.User.String())
		}

		return c.EditC(minidis.EditProps{
			Content: "I have successfully removed the `Linked Pup` role of the following members/users in the server. ",
			Attachments: []*discordgo.File{
				{
					ContentType: "text/plain",
					Name:        "unlinked-users-members.txt",
					Reader:      strings.NewReader(strings.Join(names, "\n")),
				},
			},
		})

	},
}
