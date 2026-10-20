package main

import (
	"os"

	"sylr.dev/fix/v2/cmd"
)

func main() {
	err := cmd.FixCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
