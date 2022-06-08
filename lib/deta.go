package lib

import (
	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key    string `json:"key"`
	ID     string `json:"id"`
	Wallet string `json:"wallet"`
}

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
