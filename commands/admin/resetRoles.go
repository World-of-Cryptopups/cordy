package admin

import (
	"fmt"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

func CheckPermission() {

}

func GetAllMembers(guildId string, s *discordgo.Session) []*discordgo.Member {
	lastId := "0"
	allMembers := []*discordgo.Member{}

	filteredMembers := []*discordgo.Member{}

	for {
		mems, err := s.GuildMembers(guildId, lastId, 1000)
		if err != nil {
			lib.LogError(err)
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
	Name:        "reset-roles",
	Description: "Reset the roles of all members in the server",
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(true)

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

		allMembers := GetAllMembers(c.GuildId, c.Session)

		for _, v := range allMembers {
			// check and remove the roles of all of the members in here
			fmt.Println(v)
		}

		return c.Edit("members length ")
	},
}
