package commands

import (
	"fmt"
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/atomicassets-go"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

func GetAccountStats(wallet string) atomicassets.AccountCollectionDataProps {
	q, err := lib.Atom.GetAccountCollection(wallet, "cryptopuppie")
	lib.LogError(err)

	return q.Data
}

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

		stats := GetAccountStats(user.Wallet)

		embedFields := []*discordgo.MessageEmbedField{
			{
				Name:  ":credit_card: Wallet",
				Value: user.Wallet,
			},
			{
				Name:  "\u200b",
				Value: "~ your current collection stats ~",
			},
		}
		for _, v := range stats.Schemas {
			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   v.SchemaName,
				Value:  v.Assets,
				Inline: true,
			})
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
			Fields:      embedFields,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		return c.EditC(minidis.EditProps{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
	},
}
