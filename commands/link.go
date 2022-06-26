package commands

import (
	"fmt"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/lib"
	"github.com/World-of-Cryptopups/cordy/lib/dps"
	"github.com/bwmarrin/discordgo"
	"github.com/deta/deta-go/service/base"
)

var linkCommand = &minidis.SlashCommandProps{
	Name:        "link",
	Description: "Link your wallet with your discord id",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "token",
			Description: "Your authentication token",
			Required:    true,
		},
	},
	Execute: func(c *minidis.SlashContext) error {
		c.DeferReply(true)

		token := c.Options["token"].StringValue()
		if token == "" {
			return c.Edit("token cannot be empty")
		}

		userid := c.Author.ID

		_, exists := lib.GetUser(userid)
		if exists {
			// if already linked, show error
			return c.Edit("You currently have a linked wallet. If you want to change, please contact an admin. Thank you~")
		}

		login := lib.FetchWebLoginToken(token)
		if login.Token == "" {
			return c.Edit("This token does not exist, please get your own authentication token at https://worldofcryptopups.cf/my-collections")
		}

		// check if token has been used already
		if login.Linked {
			return c.Edit("Token has already been used. If there was a mistake, please contact an admin. Thank you~")
		}

		usersBase := lib.UsersBase()
		loginBase := lib.WebLoginBase()

		newUser := lib.User{
			Key:    userid,
			ID:     userid,
			Wallet: login.Wallet,
			Token:  token,
		}

		// add to users base
		if err := loginBase.Update(login.Wallet, base.Updates{
			"linked": true,
		}); err != nil {
			_, e := c.Followup("Failed to update login, please report this error to admin.")
			return e
		}

		// update web login info
		if _, err := usersBase.Put(newUser); err != nil {
			lib.PrintError(err)

			_, err = c.Followup("There was a problem trying to link your wallet, if the problem persists, please contact an admin.")
			return err
		}

		data := dps.Calculate(login.Wallet)
		totalDps := data.PuppyCards + data.PupSkinCards + data.PupItems.Real

		// add `Verified Pups` role
		if err := c.Session.GuildMemberRoleAdd(c.GuildId, userid, lib.VERIFIED_ROLE); err != nil {
			fmt.Println(err)

			_, err := c.Followup("There was a problem trying to promote the user. If the problem persists please contact an admin.")
			return err
		}

		// add pup roles
		if err := lib.HandleRole(c.Session, userid, c.GuildId, totalDps); err != nil {
			fmt.Println(err)

			_, err := c.Followup("There was a problem trying to promote the user. If the problem persists please contact an admin.")
			return err
		}

		return c.Edit("Successfully linked your wallet with your User ID. You can now check you DPS stats with `/dps` command and your role will be updated.")
	},
}
