// Package optionstest sets, in tests, the flags of commands.
package optionstest

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Set resets the flags of cmd to their defaults, then sets flags, name and
// value pairs, as given on the command line. The flags are reset again at
// the end of t.
func Set(t testing.TB, cmd *cobra.Command, flags ...string) {
	t.Helper()

	reset(t, cmd)
	t.Cleanup(func() { reset(t, cmd) })

	if len(flags)%2 != 0 {
		t.Fatalf("odd number of flag names and values: %q", flags)
	}

	for i := 0; i < len(flags); i += 2 {
		if err := cmd.Flags().Set(flags[i], flags[i+1]); err != nil {
			t.Fatalf("--%s %q: %v", flags[i], flags[i+1], err)
		}
	}
}

// Run runs cmd's PersistentPreRunE, then its RunE, with ctx, as cobra does,
// and returns their error and what cmd wrote out.
func Run(ctx context.Context, cmd *cobra.Command) (string, error) {
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetContext(ctx)
	defer cmd.SetOut(nil)

	if err := cmd.PersistentPreRunE(cmd, nil); err != nil {
		return out.String(), err
	}

	err := cmd.RunE(cmd, nil)

	return out.String(), err
}

func reset(t testing.TB, cmd *cobra.Command) {
	t.Helper()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		var err error
		if slice, ok := f.Value.(pflag.SliceValue); ok {
			err = slice.Replace(nil)
		} else {
			err = f.Value.Set(f.DefValue)
		}
		if err != nil {
			t.Fatalf("resetting --%s: %v", f.Name, err)
		}
		f.Changed = false
	})
}
