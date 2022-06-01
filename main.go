package main

import (
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"sylr.dev/fix/cmd"
)

func main() {
	err := cmd.FixCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
