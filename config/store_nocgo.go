//go:build !cgo

package config

import (
	"fmt"

	"github.com/sylr/fix2/fix/session/store"

	"sylr.dev/fix/v2/errors"
)

// openSQLite fails: the SQLite store needs cgo.
func openSQLite(path, key string) (store.Store, error) {
	return nil, fmt.Errorf("%w: the sqlite3 store requires fix built with cgo", errors.Config)
}
