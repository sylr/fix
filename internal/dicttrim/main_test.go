package main

import (
	"bytes"
	"encoding/xml"
	"slices"
	"strings"
	"testing"
)

const dictionary = `<fix type="FIX" major="5" minor="0" servicepack="2">
  <header/>
  <trailer/>
  <messages>
    <message name="NewOrderSingle" msgtype="D" msgcat="app">
      <field name="ClOrdID" required="Y"/>
      <component name="Instrument" required="Y"/>
      <component name="LegGrp" required="N"/>
      <field name="Text" required="N"/>
      <field name="EncodedText" required="N"/>
      <field name="EncodedTextLen" required="N"/>
    </message>
    <message name="NewOrderMultileg" msgtype="AB" msgcat="app">
      <field name="ClOrdID" required="Y"/>
      <component name="LegGrp" required="Y"/>
      <component name="Multileg" required="Y"/>
    </message>
  </messages>
  <components>
    <component name="Instrument">
      <field name="Symbol" required="N"/>
    </component>
    <component name="LegGrp">
      <group name="NoLegs" required="N">
        <field name="LegSymbol" required="N"/>
        <field name="EncodedLegIssuerLen" required="N"/>
        <field name="EncodedLegIssuer" required="N"/>
      </group>
    </component>
    <component name="Multileg">
      <component name="LegGrp" required="Y"/>
      <field name="MultilegModel" required="N"/>
    </component>
  </components>
  <fields>
    <field number="11" name="ClOrdID" type="STRING"/>
    <field number="35" name="MsgType" type="STRING">
      <value enum="D" description="ORDER_SINGLE"/>
      <value enum="AB" description="NEW_ORDER_MULTILEG"/>
    </field>
    <field number="55" name="Symbol" type="STRING"/>
    <field number="58" name="Text" type="STRING"/>
    <field number="354" name="EncodedTextLen" type="LENGTH"/>
    <field number="355" name="EncodedText" type="DATA"/>
    <field number="555" name="NoLegs" type="NUMINGROUP"/>
    <field number="600" name="LegSymbol" type="STRING"/>
    <field number="618" name="EncodedLegIssuerLen" type="LENGTH"/>
    <field number="619" name="EncodedLegIssuer" type="DATA"/>
    <field number="1377" name="MultilegModel" type="INT"/>
  </fields>
</fix>`

func parseString(t *testing.T, s string) *node {
	t.Helper()

	root, err := parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}

	return root
}

// named returns a node of kind with a name attribute.
func named(kind, name string) *node {
	return &node{name: kind, attrs: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: name}}}
}

// names returns the name attributes of nodes.
func names(nodes []*node) []string {
	var s []string
	for _, n := range nodes {
		s = append(s, n.attr("name"))
	}
	return s
}

func TestTrim(t *testing.T) {
	root := parseString(t, dictionary)

	// LegGrp is required by NewOrderMultileg and Multileg, which only it
	// uses: discarded first, they do not prevent dropping it.
	if err := trim(root, []string{"D"}, []string{"LegGrp"}); err != nil {
		t.Fatal(err)
	}

	if got := names(root.child("messages").children); !slices.Equal(got, []string{"NewOrderSingle"}) {
		t.Errorf("messages = %v, want NewOrderSingle", got)
	}
	if got := names(root.child("components").children); !slices.Equal(got, []string{"Instrument"}) {
		t.Errorf("components = %v, want Instrument", got)
	}

	order := root.child("messages").children[0]
	if got := names(order.children); !slices.Equal(got, []string{"ClOrdID", "Instrument", "Text", "EncodedTextLen", "EncodedText"}) {
		t.Errorf("NewOrderSingle = %v, want LegGrp dropped, EncodedTextLen first", got)
	}

	// The DATA and LENGTH fields only LegGrp uses are kept, the others it
	// and Multileg use are not.
	want := []string{"ClOrdID", "MsgType", "Symbol", "Text", "EncodedTextLen", "EncodedText", "EncodedLegIssuerLen", "EncodedLegIssuer"}
	if got := names(root.child("fields").children); !slices.Equal(got, want) {
		t.Errorf("fields = %v, want %v", got, want)
	}

	// MsgType keeps its whole enumeration.
	if n := len(root.child("fields").children[1].children); n != 2 {
		t.Errorf("MsgType values = %d, want 2", n)
	}
}

func TestTrimErrors(t *testing.T) {
	for _, tc := range []struct {
		name          string
		keep, dropped []string
		want          string
	}{
		{"unknown MsgType", []string{"D", "ZZ"}, nil, `no message with MsgType "ZZ"`},
		{"unknown component", []string{"D"}, []string{"Nope"}, `unknown component "Nope"`},
		{"required component", []string{"AB"}, []string{"LegGrp"}, `component "LegGrp" is required`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := trim(parseString(t, dictionary), tc.keep, tc.dropped)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("trim() = %v, want %s", err, tc.want)
			}
		})
	}
}

func TestCollect(t *testing.T) {
	root := parseString(t, dictionary)

	components := map[string]*node{}
	for _, c := range root.child("components").children {
		components[c.attr("name")] = c
	}

	multileg := root.child("messages").children[1]

	fields, used := map[string]bool{}, map[string]bool{}
	if err := collect(multileg.children, components, nil, fields, used); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"ClOrdID", "NoLegs", "LegSymbol", "EncodedLegIssuer", "MultilegModel"} {
		if !fields[f] {
			t.Errorf("field %s not collected", f)
		}
	}
	if !used["LegGrp"] || !used["Multileg"] || used["Instrument"] {
		t.Errorf("components = %v, want LegGrp and Multileg", used)
	}

	// Dropped components are not followed.
	fields, used = map[string]bool{}, map[string]bool{}
	if err := collect(multileg.children, components, []string{"LegGrp"}, fields, used); err != nil {
		t.Fatal(err)
	}
	if used["LegGrp"] || fields["LegSymbol"] || !fields["MultilegModel"] {
		t.Errorf("fields = %v, components = %v, want LegGrp not followed", fields, used)
	}

	if err := collect([]*node{named("component", "Nope")}, components, nil, fields, used); err == nil {
		t.Errorf("collect() of an unknown component = nil, want an error")
	}
}

func TestDrop(t *testing.T) {
	root := parseString(t, dictionary)
	order := root.child("messages").children[0]

	children, err := drop(order.children, []string{"LegGrp"})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(children); slices.Contains(got, "LegGrp") || len(got) != len(order.children)-1 {
		t.Errorf("drop() = %v, want LegGrp alone dropped", got)
	}

	// Nested references are dropped too, a required one refused.
	multileg := root.child("components").children[2]
	if _, err := drop([]*node{multileg}, []string{"LegGrp"}); err == nil || !strings.Contains(err.Error(), "component Multileg") {
		t.Errorf("drop() = %v, want LegGrp required by Multileg", err)
	}
}

func TestLengthFirst(t *testing.T) {
	types := map[string]string{
		"Text":                "STRING",
		"EncodedText":         "DATA",
		"EncodedTextLen":      "LENGTH",
		"EncodedLegIssuer":    "DATA",
		"EncodedLegIssuerLen": "LENGTH",
		"XmlData":             "XMLDATA",
		"XmlDataLen":          "LENGTH",
		"Signature":           "DATA",
		"SignatureLength":     "LENGTH",
	}
	field := func(name string) *node { return named("field", name) }

	group := named("group", "NoLegs")
	group.children = []*node{field("EncodedLegIssuer"), field("EncodedLegIssuerLen")}
	items := []*node{
		field("EncodedText"), field("Text"), field("EncodedTextLen"),
		field("XmlDataLen"), field("XmlData"),
		field("Signature"), field("SignatureLength"),
		group,
	}

	got := names(lengthFirst(items, types))
	want := []string{"EncodedTextLen", "EncodedText", "Text", "XmlDataLen", "XmlData", "SignatureLength", "Signature", "NoLegs"}
	if !slices.Equal(got, want) {
		t.Errorf("lengthFirst() = %v, want %v", got, want)
	}
	if got := names(group.children); !slices.Equal(got, []string{"EncodedLegIssuerLen", "EncodedLegIssuer"}) {
		t.Errorf("group = %v, want its LENGTH field first", got)
	}
}

func TestWrite(t *testing.T) {
	root := parseString(t, `<fix><fields><field number="1" name="A" type="STRING"><value enum="&lt;" description="a &quot;b&quot; &amp; c"/></field></fields></fix>`)

	var buf bytes.Buffer
	if err := write(&buf, root, 0); err != nil {
		t.Fatal(err)
	}

	again := parseString(t, buf.String())
	value := again.child("fields").children[0].children[0]
	if value.attr("enum") != "<" || value.attr("description") != `a "b" & c` {
		t.Errorf("written values = %q, %q, want them escaped and read back", value.attr("enum"), value.attr("description"))
	}
}
