package lib

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

// create new deta instance
func Deta() *deta.Deta {
	d, err := deta.New()
	if err != nil {
		log.Fatalln(err)
	}

	return d
}

// create / fetch the deta base
func UsersBase() *base.Base {
	d := Deta()

	base, err := base.New(d, "Users")
	if err != nil {
		log.Fatalln(err)
	}

	return base
}

// dps info deta base
func UsersDpsBase() *base.Base {
	d := Deta()

	base, err := base.New(d, "DpsDB")
	if err != nil {
		log.Fatalln(err)
	}

	return base
}

// base on the web login
func WebLoginBase() *base.Base {
	d := Deta()

	base, err := base.New(d, "WebLogin")
	if err != nil {
		log.Fatalln(err)
	}

	return base
}

// gets the wallet with the token
func FetchWebLoginToken(token string) WebLoginUserProps {
	var login []WebLoginUserProps
	loginsBase := WebLoginBase()

	_, err := loginsBase.Fetch(
		&base.FetchInput{
			Q: base.Query{
				{"token": token},
			},
			Dest:  &login,
			Limit: 1,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	if len(login) == 0 {
		return WebLoginUserProps{}
	}

	return login[0]
}

// fetch a user,
// returns a user struct and bool if the user exists
func GetUser(userid string) (User, bool) {
	usersBase := UsersBase()
	var user User

	if err := usersBase.Get(userid, &user); err != nil {
		if err == deta.ErrNotFound {
			return user, false
		}

		log.Fatalln(err)
	}

	return user, true
}

// fetch wallet
func GetUserWallet(wallet string) ([]User, error) {
	usersBase := UsersBase()

	var queryResults []User

	_, err := usersBase.Fetch(&base.FetchInput{
		Q: base.Query{
			{"wallet": wallet},
		},
		Dest:  &queryResults,
		Limit: 1,
	})

	return queryResults, err
}

// gets all of the users in db
func GetAllUser() ([]User, error) {
	var users []User

	usersBase := UsersBase()

	_, err := usersBase.Fetch(
		&base.FetchInput{
			Q:     base.Query{},
			Dest:  &users,
			Limit: 0,
		},
	)
	if err != nil {
		return []User{}, err
	}

	return users, nil
}

// update the main user's data
func UpdateUserData(userid string) error {
	usersBase := UsersBase()

	return usersBase.Update(userid, base.Updates{
		"is_stopped": false,
	})
}

// updates the user's dps
func UpdateUserDps(userid string, dps DPSProps, wallet string) error {
	dpsBase := UsersDpsBase()

	_, err := dpsBase.Put(
		UserDpsProps{
			Key:       userid,
			ID:        userid,
			Dps:       dps,
			Wallet:    wallet,
			IsStopped: false,
		},
	)

	return err
}

// get the user's dpds
func GetUserDps(userid string) (UserDpsProps, error) {
	var dps UserDpsProps
	dpsBase := UsersDpsBase()

	err := dpsBase.Get(userid, &dps)

	return dps, err
}

// update the user's data from database
func StopUser(userid string, wallet string) error {
	usersBase := UsersBase()
	dpsBase := UsersDpsBase()

	// remove id from registered users
	if err := usersBase.Update(userid, base.Updates{
		"is_stopped": true,
	}); err != nil {
		// sendlog
		SendLog(&LogProps{
			Type:        LogTypeError,
			Title:       "Failed to update user",
			Description: fmt.Sprintf("Error in updating the user's (**`%s`**) data in the database. Please update it manually.", userid),
			Message:     "Trying to update `is_stopped` key in database for user but failed.",
		})

		return fmt.Errorf("failed to update user from database. (user: %s)", userid)
	}

	// remove data from the dps database
	if err := dpsBase.Update(userid, base.Updates{
		"is_stopped": true,
	}); err != nil {
		// sendlog
		SendLog(&LogProps{
			Type:        LogTypeError,
			Title:       "Failed to update user's dps data",
			Description: fmt.Sprintf("Error in update the user's (**`%s`**) dps data in the database. Please update it manually.", userid),
			Message:     "Trying to update `is_stopped` key in database for user but failed.",
		})

		return fmt.Errorf("failed to update user's dps data from database. (user: %s)", userid)
	}

	// do not unlink token, since user's data is only stopped

	return nil
}

type UserSession struct {
	Session    *discordgo.Session
	UserExists bool
	GuildID    string
}

// remove the user from the whitelist
// remove `Linked Pup` role of user
// params:
//   - userid: user's discord id
//   - wallet: users' wallet
//   - userExists: checker if user exists or not and use
//   - reason: reason user is unlinked, useful for admin dashboard
func UnlinkUser(userid string, wallet string, userSession UserSession, reason string) error {
	usersBase := UsersBase()
	loginsBase := WebLoginBase()

	// update user's keys
	if err := usersBase.Update(userid, base.Updates{
		"is_whitelisted":         false,
		"not_whitelisted_reason": reason,
	}); err != nil {
		return fmt.Errorf("failed to remove user from database. (user: %s)", userid)
	}

	// unlink discord id from wallet
	if err := loginsBase.Update(wallet, base.Updates{
		"linked": false,
	}); err != nil {
		return fmt.Errorf("failed to unlink discord userid from wax wallet. (wallet: %s)", wallet)
	}

	// remove discord role if user exists
	// will not proceed if user left
	if userSession.UserExists {
		// session should not be null in here
		if userSession.Session == nil {
			return fmt.Errorf("dev error: session should not be null")
		}

		// remove `Linked Pup` role
		userSession.Session.GuildMemberRoleRemove(userSession.GuildID, userid, VERIFIED_ROLE)
	}

	// remove user from whitelist
	if _, err := RemoveWhitelist(wallet); err != nil {
		return fmt.Errorf("failed to remove wallet (%s) from whitelist, please remove wallet from admin dashboard", wallet)
	}

	return nil

}
