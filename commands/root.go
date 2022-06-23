package commands

import (
	"fmt"
	"log"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/commands/admin"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/auto"
	"github.com/bwmarrin/discordgo"
)

var Bot *minidis.Minidis

func init() {
	Bot = minidis.New(lib.TOKEN)

	// sync to server
	Bot.SyncToGuilds(lib.GUILD...)

	Bot.OnReady(func(s *discordgo.Session, i *discordgo.Ready) {
		log.Println("Bot is ready!")

		if !lib.DEV {
			// no need to start when in dev mode
			go auto.Start(s, lib.GUILD[0])
		}
	})

	// add commands in here
	Bot.AddCommand(linkCommand)
	Bot.AddCommand(profileCommand)
	Bot.AddCommand(dpsCommand)
	Bot.AddCommand(helpCommand)
	Bot.AddCommand(infoCommand)
	Bot.AddCommand(&minidis.SlashCommandProps{
		Name:        "hi",
		Description: "hi",
		Options: []*discordgo.ApplicationCommandOption{{
			Name:        "member",
			Description: "The member that you want to get the walelt info.",
			Required:    true,
			Type:        discordgo.ApplicationCommandOptionUser,
		}},
		Execute: func(c *minidis.SlashContext) error {
			fmt.Println(c.Options)

			return c.ReplyString("hello")
		},
	})
	Bot.AddCommand(admin.ResetRolesCommand)
	Bot.AddCommand(admin.ListUnverifiedCommand)
	Bot.AddCommand(admin.KickUnverifiedCommand)
	Bot.AddCommand(admin.GetVerifiedWalletsCommand)
	Bot.AddCommand(admin.GetRoleWalletsCommand)

	// components
	Bot.AddComponentHandler(profileBtn)
	Bot.AddComponentHandler(dpsBtn)
}
