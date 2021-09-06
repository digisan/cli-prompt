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

	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println(m["Bool"])
	fmt.Println(m["IP"])
	fmt.Println("should be <nil>:", m["_IP"])
	fmt.Println("should not be <nil>:", m["for_expert"])
}
