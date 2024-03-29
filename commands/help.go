package commands

import (
	"time"

	"github.com/TheBoringDude/minidis"
	"github.com/bwmarrin/discordgo"
)

var profileBtn = &minidis.ComponentInteractionProps{
	ID: "btn-profile",
	Execute: func(s *minidis.SlashContext, c *minidis.ComponentContext) error {
		return profileCommand.Execute(s)
	},
}

var dpsBtn = &minidis.ComponentInteractionProps{
	ID: "btn-dps",
	Execute: func(s *minidis.SlashContext, c *minidis.ComponentContext) error {
		return dpsCommand.Execute(s)
	},
}

var helpCommand = &minidis.SlashCommandProps{
	Name:                     "help",
	Description:              "Show help message for the bot",
	DefaultMemberPermissions: 1 << 31,
	Execute: func(c *minidis.SlashContext) error {
		components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Profile",
						Style:    discordgo.PrimaryButton,
						CustomID: "btn-profile",
					},
					discordgo.Button{
						Label:    "DPS",
						Style:    discordgo.SecondaryButton,
						CustomID: "btn-dps",
					},
				},
			},
		}

		embed := &discordgo.MessageEmbed{
			Title: c.Bot.Username,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    c.Bot.Username,
				IconURL: c.Bot.AvatarURL(""),
			},
			Description: "I am a helper bot for World of Cryptopups server. You can use my slash commands for doing some amazing stuff...",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "/link",
					Value: "Link your wax wallet using your auth token from the website (https://www.worldofcryptopups.com/profile)",
				},
				{
					Name:  "/profile",
					Value: "Check your profile wallet and stats in the collection",
				},
				{
					Name:  "/dps",
					Value: "Calculate and shows your DPS stats with your pup collections",
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/attachments/880103250257399859/993917469364322484/Link.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "© 2022 | World of Cryptopups",
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		return c.ReplyC(minidis.ReplyProps{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
			Components: components,
		})
	},
}
