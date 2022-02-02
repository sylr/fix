package dictgen

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

type FixVer struct {
	Ver  string
	Name string
}

func Generate(xmlpaths []string, dictdir string, t *template.Template) error {
	if len(xmlpaths) == 0 {
		return errors.New("no xml spec file")
	}

	versions := make([]FixVer, 0, len(xmlpaths))
	for _, xmlpath := range xmlpaths {
		fmt.Printf("processing %s...", xmlpath)
		g, err := newGenerator(xmlpath, t)
		if err == nil {
			err = g.generate(dictdir)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		versions = append(versions, FixVer{g.getFixVersion(), g.getFixName()})
		fmt.Fprintln(os.Stderr, "")
	}

	return generateDict(versions, dictdir, t)
}

type generator struct {
	doc *XMLDoc
	t   *template.Template
}

func newGenerator(xmlpath string, t *template.Template) (*generator, error) {
	xmlfile, err := os.Open(xmlpath)
	if err != nil {
		return nil, err
	}
	defer xmlfile.Close()

	doc := new(XMLDoc)
	decoder := xml.NewDecoder(xmlfile)
	if err = decoder.Decode(doc); err != nil {
		return nil, err
	}
	return &generator{doc, t}, nil
}

func (g *generator) getFixVersion() string {
	return g.doc.Version
}

func (g *generator) getFixName() string {
	ver := strings.Replace(g.doc.Version, ".", "", -1)
	ver = strings.Replace(ver, "+", "", -1)
	return ver
}

type tagSource struct {
	FixVer  string
	FixName string
	Fields  []*XMLField
	Enums   []*XMLEnum
}

func (g *generator) generate(outdir string) error {
	fixname := g.getFixName()

	outpath := path.Join(outdir, fixname+".generated.go")
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(os.Stderr, " -> %s", outpath)

	err = g.t.ExecuteTemplate(f, "tag.txt", tagSource{
		FixVer:  g.doc.Version,
		FixName: fixname,
		Fields:  g.doc.Fields,
		Enums:   g.doc.Enums,
	})

	return nil
}

func generateDict(versions []FixVer, outdir string, t *template.Template) error {
	outpath := path.Join(outdir, "dict.generated.go")
	fmt.Fprintf(os.Stderr, "generate %s\n", outpath)

	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.ExecuteTemplate(f, "dict.txt", versions)
}
