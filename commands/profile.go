package commands

import (
	"fmt"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

var profileCommand = &minidis.SlashCommandProps{
	Name:        "profile",
	Description: "Show your user profile",
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		userid := c.Author.ID

		user, exists := lib.GetUser(userid)
		if !exists {
			_, err := c.Followup("You haven't linked your wallet. Please use the command `/link` ")
			return err
		}

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    c.Author.Username,
				IconURL: c.Author.AvatarURL(""),
			},
			Title: fmt.Sprintf("%s - Profile", c.Author.Username),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: c.Author.AvatarURL(""),
			},
			Description: "Your user profile",
			Fields: []*discordgo.MessageEmbedField{{
				Name:  ":credit_card: Wallet",
				Value: user.Wallet,
			}},
		}

		return c.EditC(minidis.EditProps{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
	},
}
