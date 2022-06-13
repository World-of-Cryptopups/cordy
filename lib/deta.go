package lib

import (
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
