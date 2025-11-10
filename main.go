package main

import (
	"fmt"
	"os"

	"project_selector/src/cli"
)

func main() {
	err := cli.RunSelector()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
