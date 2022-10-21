package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/TheBoringDude/minidis"
	"github.com/World-of-Cryptopups/cordy/commands"
)

func main() {
	if err := minidis.Execute(commands.Bot); err != nil {
		log.Fatal(err)
	}
}
