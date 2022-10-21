package main

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/commands"
	"github.com/World-of-Cryptopups/cordy/lib"
)

func main() {
	blacklists, err := lib.GetBlacklists()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(blacklists)

	if err := minidis.Execute(commands.Bot); err != nil {
		log.Fatal(err)
	}
}
