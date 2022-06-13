package lib

import (
	"errors"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DPSStats struct {
	Title  string
	RoleID string
	Color  string
}

var initRoles = strings.Split(os.Getenv("ROLES"), ",")

const PupsWarrior = 3000
const PupsKnight = 5000
const PupsOverlord = 8000
const PupsApocalypse = 10000
const PupsAboveAll = 25000
const PupsDoggosOfInfinity = 70000
const PupsDoggosOfEternity = 145000

// Roles is the roles and
var Roles = map[int]DPSStats{
	PupsWarrior: {
		Title:  "Warrior Pups",
		RoleID: initRoles[0],
		Color:  "green",
	},
	PupsKnight: {
		Title:  "Knight Pups",
		RoleID: initRoles[1],
		Color:  "blue",
	},
	PupsOverlord: {
		Title:  "Overlord Pups",
		RoleID: initRoles[2],
		Color:  "purple",
	},
	PupsApocalypse: {
		Title:  "Pups of the Apocalypse",
		RoleID: initRoles[3],
		Color:  "red",
	},
	PupsAboveAll: {
		Title:  "Pups Above All",
		RoleID: initRoles[4],
		Color:  "orange",
	},
	PupsDoggosOfInfinity: {
		Title:  "Doggos of Infinity",
		RoleID: initRoles[5],
		Color:  "gold",
	},
	PupsDoggosOfEternity: {
		Title:  "Doggos of Eternity",
		RoleID: initRoles[6],
		Color:  "white",
	},
}

// AllRoles is the list of all available roles for ranking.
var AllRoles = []string{
	"Warrior Pups",
	"Knight Pups",
	"Overlord Pups",
	"Pups of the Apocalypse",
	"Pups Above All",
	"Doggos of Infinity",
	"Doggos of Eternity",
}

var Colors = map[string]string{
	"purple": "#a652bb",
	"blue":   "#3b6fff",
	"cyan":   "#00c09a",
	"green":  "#00d166",
	"gold":   "#fff000",
	"red":    "#e61616",
	"orange": "#ffa500",
	"white":  "#ffffff",
	"grey":   "#95a5a6",
}

// GetDPSRoleInfo gets the role info for a specific DPS.
func GetDPSRoleInfo(dps int) DPSStats {
	var d DPSStats

	if dps >= PupsWarrior && dps < PupsKnight {
		d = Roles[PupsWarrior]
	} else if dps >= PupsKnight && dps < PupsOverlord {
		d = Roles[PupsKnight]
	} else if dps >= PupsOverlord && dps < PupsApocalypse {
		d = Roles[PupsOverlord]
	} else if dps >= PupsApocalypse && dps < PupsAboveAll {
		d = Roles[PupsApocalypse]
	} else if dps >= PupsAboveAll && dps < PupsDoggosOfInfinity {
		d = Roles[PupsAboveAll]
	} else if dps >= PupsDoggosOfInfinity && dps < PupsDoggosOfEternity {
		d = Roles[PupsDoggosOfInfinity]
	} else if dps >= PupsDoggosOfEternity {
		d = Roles[PupsDoggosOfEternity]
	}

	return d
}

// promotes the user with his/her dps stats
func HandleRole(session *discordgo.Session, userid string, guildId string, dps int) error {
	d := GetDPSRoleInfo(dps)

	if d.Title == "" {
		return errors.New("there was a problem trying to promote the user")
	}

	// promote the user
	for i, v := range Roles {
		if i <= dps {
			if err := session.GuildMemberRoleAdd(guildId, userid, v.RoleID); err != nil {
				return err
			}
		} else {
			if err := session.GuildMemberRoleRemove(guildId, userid, v.RoleID); err != nil {
				return err
			}
		}
	}

	return nil
}
