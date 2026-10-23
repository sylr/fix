// Package fix44 holds the FIX 4.4 dictionary fix validates FIX.4.4
// sessions with: QuickFIX's FIX44.xml, from quickfixgo's spec directory,
// under the QuickFIX Software License.
package fix44

import (
	"bytes"
	"compress/bzip2"
	_ "embed"
	"fmt"
	"sync"

	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/dictionary/xmldict"
)

// BeginString is the BeginString (8) of FIX 4.4 messages.
const BeginString = "FIX.4.4"

//go:embed FIX44.xml.bz2
var xml []byte

// Dictionary returns the FIX 4.4 dictionary, loaded on first use and shared:
// it must not be changed.
var Dictionary = sync.OnceValue(func() *dictionary.Dict {
	d, err := xmldict.Load(bzip2.NewReader(bytes.NewReader(xml)))
	if err != nil {
		panic(fmt.Sprintf("fix44: loading the embedded dictionary: %v", err))
	}
	return d
})
