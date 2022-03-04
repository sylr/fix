//go:build acceptor
// +build acceptor

package cmd

import (
	"sylr.dev/fix/cmd/acceptor"
)

func init() {
	FixCmd.AddCommand(acceptor.AcceptorCmd)
}
