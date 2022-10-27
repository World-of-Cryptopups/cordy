package auto

import (
	"fmt"
	"log"
	"time"

	"github.com/World-of-Cryptopups/cordy/commands/admin"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/dps"
	"github.com/World-of-Cryptopups/cordy/utils"
	"github.com/bwmarrin/discordgo"
)

// start the auto bot fetcher
func Start(session *discordgo.Session, guildId string) {
	for {
		fmt.Println("- starting worker -")

		// checker to check the number of wallets that has been updated
		updatedWallets := 0

		blacklists, err := lib.GetBlacklists()
		if err != nil {
			// TODO: improve error hanlding in here
			log.Println(err)
		}

		users, err := lib.GetAllUser()
		if err != nil {
			log.Println(err)
		}

		for _, v := range users {
			user, err := session.GuildMember(guildId, v.ID)
			if err != nil {
				// user does not exist in guild / other problems
				fmt.Printf("%s | err : %v\n", v.ID, err)

				if errMsg, ok := err.(*discordgo.RESTError); ok {
					// 10013 == Unknown user
					// 10007 == Unknown member
					if errMsg.Message.Code == discordgo.ErrCodeUnknownUser || errMsg.Message.Code == discordgo.ErrCodeUnknownMember {
						fmt.Printf("User: %s has left and his data has been stopped\n", v.Wallet)

						// send log
						lib.SendLog(&lib.LogProps{
							Type:        lib.LogTypeInfo,
							Title:       "User Left",
							Description: fmt.Sprintf("User has left the server (**`%s`** - `%s`) and will be removed from the database", v.ID, v.Wallet),
							Message:     fmt.Sprintf("`%v`", err),
						})

						// remove the user from the database
						// - this is for the purpose to remove them from the /leaderboard page
						//    if they left the server
						// TODO: user should be removed if he is not in server
						if err = lib.StopUser(v.ID, v.Wallet); err != nil {
							log.Printf("Error: %v\n", err)
						}

						continue
					}
				}

				// send log
				lib.SendLog(&lib.LogProps{
					Type:        lib.LogTypeInfo,
					Title:       "User Get Failed",
					Description: fmt.Sprintf("Failed to get user: **`%s`** and has stopped to update dps, I will try again later.", v.ID),
					Message:     fmt.Sprintf("`%v`", err),
				})

				continue
			}

			// check if they are blacklisted
			if utils.Includes(v.Wallet, blacklists) {
				fmt.Printf("Wallet is blacklisted: %s | .. Removing it from db\n", v.Wallet)

				// send log
				lib.SendLog(&lib.LogProps{
					Type:        lib.LogTypeInfo,
					Title:       "User has been blacklisted from the services",
					Description: "The user's wallet is included in the blacklist system of the project",
					Message:     fmt.Sprintf("The user's wallet: **`%s`** has been blacklisted, all of his data will be removed from the database", v.Wallet),
				})

				// if wallet is blacklisted, remove from db
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

			// promote user
			if err = lib.HandleRole(session, v.ID, guildId, totalDps); err != nil {
				log.Println(err)
			}

			// update the user's data
			if err = lib.UpdateUserData(v.ID); err != nil {
				log.Println(err)
			}

			// update the user's dps
			if err = lib.UpdateUserDps(v.ID, data, v.Wallet); err != nil {
				log.Println(err)
			}

			updatedWallets += 1
			fmt.Println(v)
		}

		fmt.Printf("-------------------------------------------------------------\nAuto DPS worker is done, Total Wallets: %d\n-------------------------------------------------------------\n", updatedWallets)

		// sleep for 1 minute
		time.Sleep(time.Duration(1) * time.Minute)
	}
}
