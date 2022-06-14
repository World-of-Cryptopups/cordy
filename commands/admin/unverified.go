package admin

import (
	"fmt"
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

func hasRole(r string, roles []string) bool {
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
		if hasRole(lib.ADMIN_ROLE, v.Roles) || hasRole(lib.MOD_ROLE, v.Roles) {
			continue
		}

		if !hasRole(lib.ADVENTURE_ROLE, v.Roles) {
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

		memsStr := ""
		for _, v := range mems {
			memsStr += fmt.Sprintf("`%s`\n", v.User.Username)
		}

		embed := &discordgo.MessageEmbed{
			Title: "List of Unverified Members",
			Description: fmt.Sprintf(`The following are users/members that are verified in the server

	%s
			`, memsStr),
		}

		return c.EditC(minidis.EditProps{
			Embeds: []*discordgo.MessageEmbed{
				embed,
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

		return c.Edit("- not yet done huhu")
	},
}
