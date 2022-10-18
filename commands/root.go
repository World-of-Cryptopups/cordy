package commands

import (
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

	// admin commands in here
	Bot.AddCommand(admin.ResetRolesCommand)
	Bot.AddCommand(admin.ListUnverifiedCommand)
	Bot.AddCommand(admin.KickUnverifiedCommand)
	Bot.AddCommand(admin.GetVerifiedWalletsCommand)
	Bot.AddCommand(admin.GetRoleWalletsCommand)
	Bot.AddCommand(admin.ResetAccountCommand)
	Bot.AddCommand(admin.FindWalletCommand)

	// components
	Bot.AddComponentHandler(profileBtn)
	Bot.AddComponentHandler(dpsBtn)
}
