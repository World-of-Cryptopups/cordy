package auto

import (
	"fmt"
	"log"
	"time"

	"github.com/World-of-Cryptopups/cordy/commands/admin"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/dps"
	"github.com/bwmarrin/discordgo"
	"github.com/deta/deta-go/service/base"
	"github.com/tbdsux/mini-go/mini"
)

// start the auto bot fetcher
func Start(session *discordgo.Session, guildId string) {
	for {
		fmt.Println("- starting worker -")

		// checker to check the number of wallets that has been updated
		updatedWallets := 0

		blacklists, err := lib.GetBlacklists()
		if err != nil {
			log.Fatalf("there was a problem getting the blacklists from contract, pls contact dev\n err: %v\n", err)
			lib.SendLog(&lib.LogProps{
				Type:        lib.LogTypeError,
				Title:       "Failed to get blacklists",
				Description: "pls contact dev to check error",
				Message:     fmt.Sprintf("%v", err),
			})
		}

		whitelists, err := lib.GetWhitelists()
		if err != nil {
			log.Fatalf("there was a problem getting the whitelists from contract, pls contact dev\n err: %v\n", err)
			lib.SendLog(&lib.LogProps{
				Type:        lib.LogTypeError,
				Title:       "Failed to get whitelists",
				Description: "pls contact dev to check error",
				Message:     fmt.Sprintf("%v", err),
			})
		}

		users, err := lib.GetAllUser()
		if err != nil {
			log.Fatalf("there was a problem getting all users from database, pls contact dev\n err: %v\n", err)
			lib.SendLog(&lib.LogProps{
				Type:        lib.LogTypeError,
				Title:       "Failed to get users from db",
				Description: "pls contact dev to check error",
				Message:     fmt.Sprintf("%v", err),
			})
		}

		for _, v := range users {
			// check if user exists in the discord guild / server
			user, err := session.GuildMember(guildId, v.ID)
			if err != nil {
				// user does not exist in guild / other problems
				fmt.Printf("%s | err : %v\n", v.ID, err)

				if errMsg, ok := err.(*discordgo.RESTError); ok {
					// 10013 == Unknown user
					// 10007 == Unknown member
					if errMsg.Message.Code == discordgo.ErrCodeUnknownUser || errMsg.Message.Code == discordgo.ErrCodeUnknownMember {
						fmt.Printf("User: %s has left and will be removed from the database\n", v.Wallet)

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
						if err = lib.UnlinkUser(v.ID, v.Wallet, lib.UserSession{
							Session:    session,
							GuildID:    guildId,
							UserExists: false,
						}, "user left"); err != nil {
							log.Printf("Error: %v\n", err)
						}

						continue
					}
				}

				// send log
				lib.SendLog(&lib.LogProps{
					Type:        lib.LogTypeError,
					Title:       "User Get Failed",
					Description: fmt.Sprintf("Failed to get user: **`%s`** and has stopped to update dps, I will try again later.", v.ID),
					Message:     fmt.Sprintf("`%v`", err),
				})

				continue
			}

			// check if they are blacklisted
			if mini.Exists(blacklists, v.Wallet) {
				fmt.Printf("Wallet is blacklisted: %s | .. Removing it from db\n", v.Wallet)

				// send log
				lib.SendLog(&lib.LogProps{
					Type:        lib.LogTypeInfo,
					Title:       "User has been blacklisted from the services",
					Description: fmt.Sprintf("<!%s>'s wallet exists in the contract's blacklists, please check user.", v.ID),
					Message:     fmt.Sprintf("The user's wallet: **`%s`** has been blacklisted, users' data will be marked from db", v.Wallet),
				})

				go func() {
					// we do not want this function to block the current process
					// if wallet is blacklisted, auto-unlink user
					if err = lib.UnlinkUser(v.ID, v.Wallet, lib.UserSession{
						UserExists: true,
						Session:    session,
						GuildID:    guildId,
					}, "blacklisted"); err != nil {
						log.Printf("Error: %v\n", err)
					}
				}()

				continue
			}

			// check if wallet exists in the whitelist list from contract
			if !mini.Exists(whitelists, v.Wallet) {
				// if wallet does not exist, update user's data key `is_whitelisted` to false
				lib.SendLog(&lib.LogProps{
					Type:        lib.LogTypeInfo,
					Title:       "Missing from whitelist",
					Description: fmt.Sprintf("<!%s>'s wallet does not exist in whitelist, please re-check with user.", v.ID),
					Message:     fmt.Sprintf("The user's wallet: **`%s`** has is missing from contract's whitelist.", v.Wallet),
				})

				go func() {
					// we do not want this function to block the current process
					// if wallet is blacklisted, auto-unlink user
					if err = lib.UnlinkUser(v.ID, v.Wallet, lib.UserSession{
						UserExists: true,
						Session:    session,
						GuildID:    guildId,
					}, "missing from whitelist"); err != nil {
						log.Printf("Error: %v\n", err)
					}
				}()

				continue
			} else {
				// TODO: might have a better function than this

				usersBase := lib.UsersBase()

				if err := usersBase.Update(v.ID, base.Updates{
					"is_whitelisted":         true,
					"not_whitelisted_reason": "",
				}); err != nil {
					log.Printf("failed to update user (user is in whitelist), err: %v\n", err)
				}
			}

			// add `Linked Pups` role if the user is registered but doesn't have it
			if !admin.HasRole(lib.VERIFIED_ROLE, user.Roles) {
				session.GuildMemberRoleAdd(guildId, user.User.ID, lib.VERIFIED_ROLE)
			}

			// == calculate dps in
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
