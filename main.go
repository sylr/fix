package main

import (
	"os"

	"sylr.dev/fix/cmd"
)

func main() {
	err := cmd.FixCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
