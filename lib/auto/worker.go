package auto

import (
	"fmt"
	"log"
	"time"

	"github.com/World-of-Cryptopups/cordy/commands/admin"
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
			user, err := session.GuildMember(guildId, v.ID)
			if err != nil {
				// user does not exist in guild / other problems
				fmt.Printf("%s | err : %v\n", v.ID, err)

				// remove the user from the database
				// - this is for the purpose to remove them from the /leaderboard page
				//    if they left the server
				if err = lib.RemoveUser(v.ID, v.Wallet); err != nil {
					log.Printf("Error: %v\n", err)
				}

				continue
			}

			// add `Verified Pups` role if the user is registered but doesn't have it
			if !admin.HasRole(lib.VERIFIED_ROLE, user.Roles) {
				session.GuildMemberRoleAdd(guildId, user.User.ID, lib.VERIFIED_ROLE)
			}

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
