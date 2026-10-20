//go:build acceptor
// +build acceptor

package cmd

import (
	"sylr.dev/fix/v2/cmd/acceptor"
)

func init() {
	FixCmd.AddCommand(acceptor.AcceptorCmd)
}
