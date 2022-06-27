package admin

import (
	"fmt"
	"strings"
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

func HasRole(r string, roles []string) bool {
	for _, v := range roles {
		if r == v {
			return true
		}
	}

	return false
}

func GetAllUnverifiedMembers(guildId string, session *discordgo.Session) []*discordgo.Member {
	allMembers := GetAllMembers(guildId, session)

	unverifiedMembers := []*discordgo.Member{}

	for _, v := range allMembers {
		// pass if a mod or admin
		if HasRole(lib.ADMIN_ROLE, v.Roles) || HasRole(lib.MOD_ROLE, v.Roles) {
			continue
		}

		if !HasRole(lib.ADVENTURE_ROLE, v.Roles) {
			// do not kick if member only joined within a day
			if time.Since(v.JoinedAt).Hours() < 24 {
				continue
			}

			unverifiedMembers = append(unverifiedMembers, v)
		}
	}

	return unverifiedMembers
}

var ListUnverifiedCommand = &minidis.SlashCommandProps{
	Name:        "list-unverified",
	Description: "Lists the users in the server that are not verified",
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

		mems := GetAllUnverifiedMembers(c.GuildId, c.Session)

		strMems := []string{}
		for _, v := range mems {
			strMems = append(strMems, v.User.Username)
		}

		parsedList := strings.Join(strMems, "\n")

		return c.EditC(minidis.EditProps{
			Content: "The following are users/members that are not yet verified in the server.",
			Attachments: []*discordgo.File{
				{
					ContentType: "text/plain",
					Name:        "unverified.txt",
					Reader:      strings.NewReader(parsedList),
				},
			},
		})

	},
}

var KickUnverifiedCommand = &minidis.SlashCommandProps{
	Name:        "kick-unverified",
	Description: "Kick the users that are not verified in the server",
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

		mems := GetAllUnverifiedMembers(c.GuildId, c.Session)
		for _, v := range mems {
			if err := c.Session.GuildMemberDelete(c.GuildId, v.User.ID); err != nil {
				_, e := c.Followup(fmt.Sprintf("Failed to kick member: `%s`", v.User.Username))
				fmt.Println(e)

				return e
			}
		}

		return c.Edit("Successfully kicked all unverified members in the server.")
	},
}
