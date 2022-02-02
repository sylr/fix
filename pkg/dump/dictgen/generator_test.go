package dictgen_test

import (
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"

	"sylr.dev/fix/pkg/dump/dictgen"
)

func TestGenerator(t *testing.T) {
	xmlpaths, err := filepath.Glob("xml/*.xml")
	require.NoError(t, err)

	tmpl, err := template.ParseGlob("tmpl/*.txt")
	require.NoError(t, err)

	err = dictgen.Generate(xmlpaths, "../dict", tmpl)
	require.NoError(t, err)
}
