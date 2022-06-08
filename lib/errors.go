package lib

import (
	"fmt"
	"log"
)

func LogError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
