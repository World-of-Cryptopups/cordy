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

		verifiedMems := []*discordgo.Member{}
		users := GetAllMembers(c.GuildId, c.Session)

		for _, v := range users {
			if hasRole(lib.VERIFIED_ROLE, v.Roles) {
				verifiedMems = append(verifiedMems, v)
			}
		}

		wallets := []string{}
		for _, v := range verifiedMems {
			user, exists := lib.GetUser(v.User.ID)
			if !exists {
				continue
			}

			wallets = append(wallets, user.Wallet)
		}

		return c.EditC(minidis.EditProps{
			Content: fmt.Sprintf("Here are the wallets of the verified members / users...\n\nTotal members: %d \t", len(wallets)),
			Attachments: []*discordgo.File{{
				ContentType: "text/plain",
				Name:        "wallets.txt",
				Reader:      strings.NewReader(strings.Join(wallets, "\n")),
			}},
		})
	},
}

var GetRoleWalletsCommand = &minidis.SlashCommandProps{
	Name:        "get-role-wallets",
	Description: "Get wallets of members that have role.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "role",
			Description: "Role of the members that you want to get",
			Required:    true,
			Type:        discordgo.ApplicationCommandOptionRole,
		},
	},
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		role := c.Options["role"].RoleValue(c.Session, c.GuildId)

		roleMems := []*discordgo.Member{}
		mems := GetAllMembers(c.GuildId, c.Session)

		for _, v := range mems {
			if hasRole(role.ID, v.Roles) {
				roleMems = append(roleMems, v)
			}
		}

		wallets := []string{}
		for _, v := range roleMems {
			user, exists := lib.GetUser(v.User.ID)
			if !exists {
				continue
			}

			wallets = append(wallets, user.Wallet)
		}

		return c.EditC(minidis.EditProps{
			Content: fmt.Sprintf("Here are the wallets of the members that have role: **%s**", role.Name),
			Attachments: []*discordgo.File{
				{
					ContentType: "text/plain",
					Name:        fmt.Sprintf("%s-wallets.txt", role.Name),
					Reader:      strings.NewReader(strings.Join(wallets, "\n")),
				},
			},
		})

	},
}
