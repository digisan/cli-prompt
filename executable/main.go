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

	fmt.Println("")

	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println("")

	fmt.Println(m["Bool"])
	fmt.Println(m["IP"])
	fmt.Println(`should be <nil> @m["_IP"]:`, m["_IP"])
	fmt.Println(`should NOT be <nil> @m["for_expert"]:`, m["for_expert"])
}
