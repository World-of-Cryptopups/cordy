package admin

import (
	"fmt"
	"log"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

func GetAllMembers(guildId string, s *discordgo.Session) []*discordgo.Member {
	lastId := "0"
	allMembers := []*discordgo.Member{}

	filteredMembers := []*discordgo.Member{}

	for {
		mems, err := s.GuildMembers(guildId, lastId, 1000)
		if err != nil {
			fmt.Println("Err! Failed to get all members in server")
			log.Fatalln(err)
		}

		if len(mems) == 0 {
			break
		}

		lastId = mems[len(mems)-1].User.ID
		allMembers = append(allMembers, mems...)
	}

	// filter the members in here
	for _, v := range allMembers {
		if !v.User.Bot {
			filteredMembers = append(filteredMembers, v)
		}
	}

	return filteredMembers
}

var ResetRolesCommand = &minidis.SlashCommandProps{
	Name:                     "reset-roles",
	Description:              "Reset the roles of all members in the server",
	DefaultMemberPermissions: 0,
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		allMembers := GetAllMembers(c.GuildId, c.Session)

		for _, v := range allMembers {
			// remove role only if has adventure role (meaning it also has other roles)
			if !HasRole(lib.ADVENTURE_ROLE, v.Roles) {
				continue
			}

			fmt.Printf("removing the role of %s \n", v.User.Username)

			// remove all of the roles in here
			for _, r := range lib.ALL_ROLES {
				if err := c.Session.GuildMemberRoleRemove(c.GuildId, v.User.ID, r); err != nil {
					fmt.Printf("Failed to remove the role of %s, error: %v\n", v.User.Username, err)

					// _, e := c.Followup(fmt.Sprintf("Failed to reset the role of %s", v.User.Username))
					// return e
				}
			}
		}

		return c.Edit("Successfully reset the roles of the all of the members in the server.")
	},
}
