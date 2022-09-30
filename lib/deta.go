package lib

import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

// create new deta instance
func Deta() *deta.Deta {
	d, err := deta.New()
	LogError(err)

	return d
}

// create / fetch the deta base
func UsersBase() *base.Base {
	d := Deta()

	base, err := base.New(d, "Users")
	LogError(err)

	return base
}

// dps info deta base
func UsersDpsBase() *base.Base {
	d := Deta()

	base, err := base.New(d, "DpsDB")
	LogError(err)

	return base
}

// base on the web login
func WebLoginBase() *base.Base {
	d := Deta()

	base, err := base.New(d, "WebLogin")
	LogError(err)

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
	LogError(err)

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

		LogError(err)
	}

	return user, true
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

// updates the user's dps
func UpdateUserDps(userid string, dps DPSProps) error {
	dpsBase := UsersDpsBase()

	_, err := dpsBase.Put(
		UserDpsProps{
			Key: userid,
			ID:  userid,
			Dps: dps,
		},
	)

	return err
}

// removes the user from the databases
func RemoveUser(userid string, wallet string) error {
	usersBase := UsersBase()
	dpsBase := UsersDpsBase()
	loginsBase := WebLoginBase()

	// remove id from registered users
	if err := usersBase.Delete(userid); err != nil {
		return fmt.Errorf("failed to remove user from database. (user: %s)", userid)
	}

	// remove data from the dps database
	if err := dpsBase.Delete(userid); err != nil {
		return fmt.Errorf("failed remove user's dps data from database. (user: %s)", userid)
	}

	// unlink discord id from wallet
	if err := loginsBase.Update(wallet, base.Updates{
		"linked": false,
	}); err != nil {
		return fmt.Errorf("failed to unlink discord userid from wax wallet. (wallet: %s)", wallet)
	}

	return nil

}
