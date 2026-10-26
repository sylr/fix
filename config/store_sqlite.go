//go:build cgo

package config

import (
	"github.com/sylr/fix2/fix/session/store"

	sqlite3store "github.com/sylr/fix2/contribs/store/sqlite3"
)

// SQLiteStore reports whether fix is built with the sqlite3 store, which
// requires cgo.
const SQLiteStore = true

// openSQLite opens the store of the session key in the SQLite database path.
func openSQLite(path, key string) (store.Store, error) {
	return sqlite3store.OpenFile(path, key)
}
