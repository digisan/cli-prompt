package main

import (
	"fmt"
	"log"

	cliprompt "github.com/digisan/cli-prompt"
)

func main() {
	m, err := cliprompt.PromptConfig("../config.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(m["Bool"])
	fmt.Println(m["IP"])
	fmt.Println(m["_IP"])
}
