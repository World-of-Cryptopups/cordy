package commands

import (
	"log"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/commands/admin"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/bwmarrin/discordgo"
)

var Bot *minidis.Minidis

func init() {
	Bot = minidis.New(lib.TOKEN)

	// sync to server
	Bot.SyncToGuilds(lib.GUILD...)

	Bot.OnReady(func(s *discordgo.Session, i *discordgo.Ready) {
		log.Println("Bot is ready!")
	})

	// add commands in here
	Bot.AddCommand(linkCommand)
	Bot.AddCommand(profileCommand)
	Bot.AddCommand(dpsCommand)
	Bot.AddCommand(admin.ResetRolesCommand)
}
