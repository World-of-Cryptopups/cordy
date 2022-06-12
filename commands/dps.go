package commands

import (
	"fmt"
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/dps"
	"github.com/bwmarrin/discordgo"
)

var dpsCommand = &minidis.SlashCommandProps{
	Name:        "dps",
	Description: "Check your pup collection dps stats",
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(false)

		userid := c.Author.ID

		user, exists := lib.GetUser(userid)
		if !exists {
			return c.Edit("You currently do not have a linked wallet to your acc.")
		}

		data := dps.Calculate(user.Wallet)
		totalDps := data.PuppyCards + data.PupSkinCards + data.PupItems.Real

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("DPS Stats - %s", c.Author.Username),
			Description: "Your collection's dps stats",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    c.Author.Username,
				IconURL: c.Author.AvatarURL(""),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: c.Author.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎴 Puppy Cards",
					Value:  fmt.Sprint(data.PuppyCards),
					Inline: true,
				},
				{
					Name:   "🃏 Pup Skins",
					Value:  fmt.Sprint(data.PupSkinCards),
					Inline: true,
				},
				{
					Name:   "⚔️ Pup Items (Raw)",
					Value:  fmt.Sprint(data.PupItems.Raw),
					Inline: true,
				},
				{
					Name:   "⚔️ Pup Items (Real)",
					Value:  fmt.Sprint(data.PupItems.Real),
					Inline: true,
				},
				{
					Name:  "\u200b",
					Value: "\u200b",
				},
				{
					Name:  "🛡 Total DPS",
					Value: fmt.Sprintf("**%d**", totalDps)},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		return c.EditC(minidis.EditProps{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		})
	},
}
