package commands

import (
	"fmt"
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

var infoCommand = &minidis.SlashCommandProps{
	Name:        "info",
	Description: "Get wallet info of a member.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "member",
			Description: "The member that you want to get the walelt info.",
			Required:    true,
			Type:        discordgo.ApplicationCommandOptionUser,
		},
	},
	DefaultMemberPermissions: 1 << 31,
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		mentioned := c.Options["member"].UserValue(c.Session)

		user, exists := lib.GetUser(mentioned.ID)
		if !exists {
			return c.Edit("The user is currently not registered to me.")
		}

		embedFields := []*discordgo.MessageEmbedField{
			{
				Name:  ":credit_card: Wallet",
				Value: user.Wallet,
			},
		}

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    mentioned.Username,
				IconURL: mentioned.AvatarURL(""),
			},
			Title: fmt.Sprintf("%s's - Profile", mentioned.Username),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: mentioned.AvatarURL(""),
			},
			Description: "User's profile",
			Fields:      embedFields,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "2022 | World of Cryptopups",
			},
		}

		return c.EditC(minidis.EditProps{
			Embeds: []*discordgo.MessageEmbed{embed},
		})

	},
}
