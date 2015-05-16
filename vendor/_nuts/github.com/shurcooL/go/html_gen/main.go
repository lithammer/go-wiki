// Package html_gen contains helper funcs for generating HTML nodes and rendering them, safe against code injection.
// Context-aware escaping is done just like in html/template.
package html_gen

import (
	"bytes"
	"html/template"

	"golang.org/x/net/html"
)

// Text returns a plain text node.
func Text(s string) *html.Node {
	return &html.Node{
		Type: html.TextNode, Data: s,
	}
}

// A returns an anchor element <a href="{{.href}}">{{.s}}</a>.
func A(s string, href template.URL) *html.Node {
	n := &html.Node{
		Type: html.ElementNode, Data: "a",
		Attr: []html.Attribute{{Key: "href", Val: string(href)}},
	}
	n.AppendChild(Text(s))
	return n
}

// RenderNodes renders a list of HTML nodes.
// Context-aware escaping is done just like in html/template when rendering nodes.
func RenderNodes(nodes ...*html.Node) (template.HTML, error) {
	var buf bytes.Buffer
	for _, node := range nodes {
		err := html.Render(&buf, node)
		if err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

// Must is a helper that wraps a call to a function returning (template.HTML, error)
// and panics if the error is non-nil.
func Must(html template.HTML, err error) template.HTML {
	if err != nil {
		panic(err)
	}
	return html
}
