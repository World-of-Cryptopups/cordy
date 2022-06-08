package main

import (
	"log"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/commands"
)

func main() {
	if err := minidis.Execute(commands.Bot); err != nil {
		log.Fatal(err)
	}
}
