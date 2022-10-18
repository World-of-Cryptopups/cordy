package admin

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/utils"
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
			if HasRole(lib.VERIFIED_ROLE, v.Roles) {
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

		role := c.Options["role"].RoleValue(c.Session, c.GuildId)

		roleMems := []*discordgo.Member{}
		mems := GetAllMembers(c.GuildId, c.Session)

		for _, v := range mems {
			if HasRole(role.ID, v.Roles) {
				roleMems = append(roleMems, v)
			}
		}

		wallets := []string{}
		for _, v := range roleMems {
			user, exists := lib.GetUser(v.User.ID)
			if !exists {
				continue
			}

			// filter each member with their dps
			data, err := lib.GetUserDps(v.User.ID)
			if err != nil {
				log.Println(err) // observe some errors for now

				// TODO: handle error in here
				continue
			}

			// check if role is in the dps roles
			if utils.Includes(role.ID, lib.InitRoles) {
				totalDps := data.Dps.PuppyCards + data.Dps.PupSkinCards + data.Dps.PupItems.Real
				userRole := lib.GetDPSRoleInfo(totalDps)

				if userRole.RoleID != role.ID {
					// if user's dps role is not the same with the role, ignore
					continue
				}
			}

			// add wallet to the list of wallets to show
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

var FindWalletCommand = &minidis.SlashCommandProps{
	Name:        "find-wallet",
	Description: "Find the user with the linked wallet.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "wallet",
			Description: "Wax wallet to find",
			Required:    true,
			Type:        discordgo.ApplicationCommandOptionString,
		},
	},
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

		wallet := c.Options["wallet"].StringValue()

		query, err := lib.GetUserWallet(wallet)
		if err != nil {
			return c.Edit(fmt.Sprintf("There was a problem trying to find the wallet. Error: %v", err))
		}

		if len(query) == 0 {
			// no users found with wallet
			return c.Edit(fmt.Sprintf("No users found with wallet: **%s**", wallet))
		}

		// user exists
		user := query[0]
		discordUser, err := c.Session.GuildMember(c.GuildId, user.ID)
		if err != nil {
			return c.Edit(fmt.Sprintf("User with wallet found, ID: **%s**. But I cannot get his info in the server.", user.ID))
		}

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    discordUser.User.Username,
				IconURL: discordUser.AvatarURL(""),
			},
			Title: fmt.Sprintf("%s | %s", user.Wallet, discordUser.User.String()),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: discordUser.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  ":bust_in_silhouette: Username",
					Value: discordUser.User.String(),
				},
				{
					Name:   ":credit_card: Wallet",
					Value:  user.Wallet,
					Inline: true,
				},
				{
					Name:   ":id: Discord ID",
					Value:  user.ID,
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "2022 | World of Cryptopups",
			},
		}

		return c.EditC(minidis.EditProps{
			Content: "User found:",
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		})
	},
}
