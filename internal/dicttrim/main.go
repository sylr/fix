// Command dicttrim trims a QuickFIX XML dictionary down to the messages given,
// the components they reference and the fields both use.
//
// The full QuickFIX FIX 5.0 SP2 dictionary expands to more items than fix2's
// xmldict accepts: fix only ships the messages it sends and receives.
package main

import (
	"bytes"
	"compress/bzip2"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type node struct {
	name     string
	attrs    []xml.Attr
	children []*node
}

func (n *node) attr(name string) string {
	for _, a := range n.attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

func (n *node) child(name string) *node {
	for _, c := range n.children {
		if c.name == name {
			return c
		}
	}
	return nil
}

func parse(r io.Reader) (*node, error) {
	dec := xml.NewDecoder(r)
	var stack []*node
	var root *node
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			n := &node{name: t.Name.Local, attrs: t.Attr}
			if len(stack) == 0 {
				root = n
			} else {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, n)
			}
			stack = append(stack, n)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		}
	}
	if root == nil {
		return nil, fmt.Errorf("empty dictionary")
	}
	return root, nil
}

func write(w *bytes.Buffer, n *node, depth int) error {
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(w, "%s<%s", indent, n.name)
	for _, a := range n.attrs {
		fmt.Fprintf(w, ` %s="`, a.Name.Local)
		if err := xml.EscapeText(w, []byte(a.Value)); err != nil {
			return err
		}
		w.WriteString(`"`)
	}
	if len(n.children) == 0 {
		w.WriteString("/>\n")
		return nil
	}
	w.WriteString(">\n")
	for _, c := range n.children {
		if err := write(w, c, depth+1); err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "%s</%s>\n", indent, n.name)
	return nil
}

// collect records the fields and components items reference, recursively,
// but the components named in dropped.
func collect(items []*node, components map[string]*node, dropped []string, fields, used map[string]bool) error {
	for _, it := range items {
		switch it.name {
		case "field":
			fields[it.attr("name")] = true
		case "group":
			fields[it.attr("name")] = true
			if err := collect(it.children, components, dropped, fields, used); err != nil {
				return err
			}
		case "component":
			name := it.attr("name")
			if used[name] || slices.Contains(dropped, name) {
				continue
			}
			c, ok := components[name]
			if !ok {
				return fmt.Errorf("unknown component %q", name)
			}
			used[name] = true
			if err := collect(c.children, components, dropped, fields, used); err != nil {
				return err
			}
		}
	}
	return nil
}

// drop removes from items, recursively, the references to the components
// named in dropped. It refuses to drop a required one.
func drop(items []*node, dropped []string) ([]*node, error) {
	var kept []*node
	for _, it := range items {
		if it.name == "component" && slices.Contains(dropped, it.attr("name")) {
			if it.attr("required") == "Y" {
				return nil, fmt.Errorf("component %q is required", it.attr("name"))
			}
			continue
		}
		if it.name == "group" || it.name == "component" || it.name == "message" {
			children, err := drop(it.children, dropped)
			if err != nil {
				return nil, fmt.Errorf("%s %s: %w", it.name, it.attr("name"), err)
			}
			it.children = children
		}
		kept = append(kept, it)
	}
	return kept, nil
}

// lengthFirst moves the LENGTH field of each DATA field found after it in
// items just before it, recursively, as FIX requires: QuickFIX's FIX 5.0 SP2
// dictionary has some the other way round, which fix2gen refuses.
func lengthFirst(items []*node, types map[string]string) []*node {
	for _, it := range items {
		if it.name == "group" || it.name == "component" {
			it.children = lengthFirst(it.children, types)
		}
	}
	for i := 0; i < len(items); i++ {
		name := items[i].attr("name")
		if items[i].name != "field" || (types[name] != "DATA" && types[name] != "XMLDATA") {
			continue
		}
		for j := i + 1; j < len(items); j++ {
			other := items[j].attr("name")
			if items[j].name != "field" || types[other] != "LENGTH" || (other != name+"Length" && other != name+"Len") {
				continue
			}
			length := items[j]
			items = slices.Delete(items, j, j+1)
			items = slices.Insert(items, i, length)
			i++
			break
		}
	}
	return items
}

func run() error {
	in := flag.String("in", "", "QuickFIX XML dictionary, bzip2 compressed if it ends with .bz2")
	out := flag.String("out", "", "trimmed dictionary")
	msgs := flag.String("msgtypes", "", "comma separated MsgType values of the messages to keep")
	dropComponents := flag.String("drop-components", "", "comma separated names of optional components to remove")
	flag.Parse()

	if *in == "" || *out == "" || *msgs == "" {
		flag.Usage()
		return fmt.Errorf("-in, -out and -msgtypes are required")
	}

	f, err := os.Open(*in)
	if err != nil {
		return err
	}

	var r io.Reader = f
	if strings.HasSuffix(*in, ".bz2") {
		r = bzip2.NewReader(f)
	}

	root, err := parse(r)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}

	var dropped []string
	if *dropComponents != "" {
		dropped = strings.Split(*dropComponents, ",")
	}

	if err := trim(root, strings.Split(*msgs, ","), dropped); err != nil {
		return err
	}

	buf := bytes.NewBufferString(xml.Header)
	if err := write(buf, root, 0); err != nil {
		return err
	}

	return os.WriteFile(*out, buf.Bytes(), 0o644)
}

// trim removes from root the messages whose MsgType keep does not list, the
// references to the components named in dropped, then the components and
// fields left unused.
func trim(root *node, keep, dropped []string) error {
	components := map[string]*node{}
	if c := root.child("components"); c != nil {
		for _, comp := range c.children {
			components[comp.attr("name")] = comp
		}
	}

	for _, name := range dropped {
		if _, ok := components[name]; !ok {
			return fmt.Errorf("unknown component %q", name)
		}
	}

	messages := root.child("messages")
	if messages == nil {
		return fmt.Errorf("no messages")
	}

	// Messages first: the components only discarded ones require may be
	// dropped.
	messages.children = slices.DeleteFunc(messages.children, func(m *node) bool {
		return !slices.Contains(keep, m.attr("msgtype"))
	})

	for _, t := range keep {
		if !slices.ContainsFunc(messages.children, func(m *node) bool { return m.attr("msgtype") == t }) {
			return fmt.Errorf("no message with MsgType %q", t)
		}
	}

	// MsgType, which no message of a FIX 5.0 application dictionary
	// references, its header being FIXT 1.1's, is kept for generated code to
	// name message types. Its enumeration, as every field's, is kept whole,
	// the discarded messages' types included: fix names any message type
	// with it, as the RefMsgType of the Rejects it receives.
	fields := map[string]bool{"MsgType": true}
	used := map[string]bool{}

	for _, section := range []string{"header", "trailer"} {
		if s := root.child(section); s != nil {
			if err := collect(s.children, components, dropped, fields, used); err != nil {
				return fmt.Errorf("%s: %w", section, err)
			}
		}
	}

	for _, m := range messages.children {
		if err := collect(m.children, components, dropped, fields, used); err != nil {
			return fmt.Errorf("message %s: %w", m.attr("name"), err)
		}
	}

	if c := root.child("components"); c != nil {
		c.children = slices.DeleteFunc(c.children, func(n *node) bool { return !used[n.attr("name")] })
	}

	if len(dropped) > 0 {
		for _, section := range []string{"header", "trailer", "messages", "components"} {
			if s := root.child(section); s != nil {
				children, err := drop(s.children, dropped)
				if err != nil {
					return err
				}
				s.children = children
			}
		}
	}

	// DATA and LENGTH fields are all kept, those only dropped components use
	// included: counterparties may still send them, and parsing a DATA field
	// without its definition would end it at an embedded SOH.
	types := map[string]string{}
	if fs := root.child("fields"); fs != nil {
		fs.children = slices.DeleteFunc(fs.children, func(n *node) bool {
			switch n.attr("type") {
			case "DATA", "XMLDATA", "LENGTH":
				return false
			}
			return !fields[n.attr("name")]
		})
		for _, f := range fs.children {
			types[f.attr("name")] = f.attr("type")
		}
	}

	for _, section := range []string{"header", "trailer", "messages", "components"} {
		if s := root.child(section); s != nil {
			for _, n := range s.children {
				n.children = lengthFirst(n.children, types)
			}
			if section == "header" || section == "trailer" {
				s.children = lengthFirst(s.children, types)
			}
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "dicttrim:", err)
		os.Exit(1)
	}
}
