package admin

import (
	"fmt"
	"strings"
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

func chunkMembersArray(members []*discordgo.Member) [][]*discordgo.Member {
	var memsList = make([][]*discordgo.Member, int(len(members)/50)+1)

	var i = 0
	for _, v := range members {
		if len(memsList) > 0 {
			if len(memsList[i]) == 50 {
				i++
			}
		}

		memsList[i] = append(memsList[i], v)
	}

	return memsList
}

func chunkEmbeds(embeds []*discordgo.MessageEmbed) [][]*discordgo.MessageEmbed {
	var grps = make([][]*discordgo.MessageEmbed, int(len(embeds)/10)+1)

	var i = 0
	for _, v := range embeds {
		if len(grps) > 0 {
			if len(grps[i]) == 10 {
				i++
			}
		}

		grps[i] = append(grps[i], v)
	}

	return grps
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

		// 	mems := chunkMembersArray(GetAllUnverifiedMembers(c.GuildId, c.Session))

		// 	embeds := []*discordgo.MessageEmbed{}
		// 	for _, x := range mems {
		// 		memsStr := ""

		// 		for _, v := range x {
		// 			memsStr += fmt.Sprintf("`%s`\n", v.User.Username)
		// 		}

		// 		embed := &discordgo.MessageEmbed{
		// 			Title: "List of Unverified Members",
		// 			Description: fmt.Sprintf(`The following are users/members that are verified in the server

		// %s
		// 			`, memsStr),
		// 		}

		// 		embeds = append(embeds, embed)
		// 	}

		// 	for i, v := range embeds {
		// 		if i == 0 {
		// 			c.EditC(minidis.EditProps{
		// 				Embeds: []*discordgo.MessageEmbed{v},
		// 			})

		// 			continue
		// 		}

		// 		c.FollowupC(minidis.FollowupProps{
		// 			Embeds: []*discordgo.MessageEmbed{v},
		// 		})
		// 	}

		// 	return nil

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
