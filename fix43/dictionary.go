// Package fix43 holds the FIX 4.3 dictionary fix validates FIX.4.3
// sessions with: QuickFIX's FIX43.xml, from quickfixgo's spec directory,
// under the QuickFIX Software License.
package fix43

import (
	"bytes"
	"compress/bzip2"
	_ "embed"
	"fmt"
	"sync"

	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/dictionary/xmldict"
)

// BeginString is the BeginString (8) of FIX 4.3 messages.
const BeginString = "FIX.4.3"

//go:embed FIX43.xml.bz2
var xml []byte

// Dictionary returns the FIX 4.3 dictionary, loaded on first use and shared:
// it must not be changed.
var Dictionary = sync.OnceValue(func() *dictionary.Dict {
	d, err := xmldict.Load(bzip2.NewReader(bytes.NewReader(xml)))
	if err != nil {
		panic(fmt.Sprintf("fix43: loading the embedded dictionary: %v", err))
	}
	return d
})
