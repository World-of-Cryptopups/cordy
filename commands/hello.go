package commands

import "github.com/TheBoringDude/minidis"

var helloCommand = &minidis.SlashCommandProps{
	Name:        "hello",
	Description: "hello command",
	Execute: func(c *minidis.SlashContext) error {
		return c.ReplyString("hello too")
	},
}
