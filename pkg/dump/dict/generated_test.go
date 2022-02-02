package dict_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sylr.dev/fix/pkg/dump/dict"
)

func TestFix44Tag(t *testing.T) {
	dict := dict.Get("FIX.4.4")

	name, ok := dict.TagName(35)
	assert.True(t, ok)
	assert.Equal(t, "MsgType", name)

	name, ok = dict.TagName(999999)
	assert.False(t, ok)
	assert.Equal(t, "", name)
}

func TestFix44Value(t *testing.T) {
	dict := dict.Get("FIX.4.4")

	name, ok := dict.ValueName(35, "A")
	assert.True(t, ok)
	assert.Equal(t, "Logon", name)

	name, ok = dict.ValueName(35, "ZZZZZ")
	assert.False(t, ok)
	assert.Equal(t, "", name)
}
