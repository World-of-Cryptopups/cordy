package commands

import (
	"log"
	"os"

	"github.com/TheBoringDude/minidis"
	"github.com/bwmarrin/discordgo"
)

var Bot *minidis.Minidis

func init() {
	Bot = minidis.New(os.Getenv("TOKEN"))

	// sync to server
	Bot.SyncToGuilds(os.Getenv("GUILD"))

	Bot.OnReady(func(s *discordgo.Session, i *discordgo.Ready) {
		log.Println("Bot is ready!")
	})

	// add commands in here
	Bot.AddCommand(helloCommand)
}
