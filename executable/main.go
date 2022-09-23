package main

import (
	"fmt"

	cliprompt "github.com/digisan/cli-prompt"
)

func main() {
	cliprompt.PromptConfig("../config.json")

	cliprompt.Show()

	fmt.Println(cliprompt.Val[string]("IP"))

	// fmt.Println(cliprompt.Val[string]("Myname"))
	// fmt.Println(cliprompt.Val[int]("Port"))

	fmt.Println(cliprompt.ValInt("Port"))
}
