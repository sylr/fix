package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	// Schedules name IANA time zones, which Windows lacks.
	_ "time/tzdata"

	"sylr.dev/fix/v2/cmd"
)

func main() {
	// The commands stop on SIGINT and SIGTERM, through cmd.Context(); a
	// second signal kills.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	context.AfterFunc(ctx, stop)

	err := cmd.FixCmd.ExecuteContext(ctx)
	stop()

	if err != nil {
		os.Exit(1)
	}
}
