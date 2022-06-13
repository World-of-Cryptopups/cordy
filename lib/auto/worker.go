package auto

import (
	"fmt"
	"log"
	"time"

	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/dps"
	"github.com/bwmarrin/discordgo"
)

// start the auto bot fetcher
func Start(session *discordgo.Session, guildId string) {
	for {
		fmt.Println("- starting worker -")

		users, err := lib.GetAllUser()
		if err != nil {
			log.Println(err)
		}

		for _, v := range users {
			data := dps.Calculate(v.Wallet)
			totalDps := data.PuppyCards + data.PupSkinCards + data.PupItems.Real

			if err = lib.HandleRole(session, v.ID, guildId, totalDps); err != nil {
				log.Println(err)
			}

			if err = lib.UpdateUserDps(v.ID, data); err != nil {
				log.Println(err)
			}

			fmt.Println(v)
		}

		// sleep for 1 minute
		time.Sleep(time.Duration(1) * time.Minute)
	}
}
