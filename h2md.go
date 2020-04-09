package h2md

import (
	"bytes"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

// H2MD H2MD struct
type H2MD struct {
	*html.Node
	ulN          int
	blockquoteN  int
	tdN          int
	tableSpliced bool
	skipNewline  bool
	replacers    map[string]CustomReplacer
}

type CustomReplacer func(val string, n *html.Node) string

// NewH2MD create H2MD with html text
func NewH2MD(htmlText string) (*H2MD, error) {
	node, err := html.Parse(strings.NewReader(htmlText))
	if err == nil {
		return &H2MD{
			Node:         node,
			ulN:          -1,
			blockquoteN:  0,
			tdN:          0,
			tableSpliced: false,
			skipNewline:  true,
			replacers:    make(map[string]CustomReplacer),
		}, nil
	}
	return nil, err
}

//NewH2MDFromNode create H2MD with html node
func NewH2MDFromNode(node *html.Node) (*H2MD, error) {
	return &H2MD{
		Node:         node,
		ulN:          -1,
		blockquoteN:  0,
		tdN:          0,
		tableSpliced: false,
		skipNewline:  true,
		replacers:    make(map[string]CustomReplacer),
	}, nil
}

// Replace Replace element attribute value
func (h *H2MD) Replace(attr string, r func(val string, n *html.Node) string) {
	h.replacers[attr] = r
}

// Attr Return the element attribute
func (h *H2MD) Attr(name string, n *html.Node) string {
	for _, attr := range n.Attr {
		if name == attr.Key {
			if r, ok := h.replacers[name]; ok {
				return r(attr.Val, n)
			}
			return attr.Val
		}
	}
	return ""
}

// Text return the markdown content
func (h *H2MD) Text() string {
	var buf bytes.Buffer

	var f func(*html.Node)

	f = func(n *html.Node) {

		var parse = func(tag string, single bool) {
			buf.WriteString(tag)
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
			if n.LastChild != nil {
				n = n.LastChild.NextSibling
			}
			if !single {
				buf.WriteString(tag)
			}
		}
		if n.Type == html.TextNode {
			if h.skipNewline {
				n.Data = strings.Trim(n.Data, "\n")
			}
			buf.WriteString(n.Data)
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "hr":
				buf.WriteString("\n---\n")
			case "a":
				if c := n.FirstChild; c != nil {
					buf.WriteString("[" + c.Data + "](" + h.Attr("href", n) + ")")
					n = c
				}
			case "img":
				if n.Parent != nil && n.Parent.Data == "\n" {
					buf.WriteString("\n")
				}
				buf.WriteString("![" + h.Attr("alt", n) + "](" + h.Attr("src", n) + ")")
			case "del":
				parse("~~", false)
			case "i":
				parse("*", false)
			case "strong", "b":
				parse("**", false)
			case "h1", "h2", "h3", "h4", "h5", "h6":
				buf.WriteString("\n")
				j, _ := strconv.Atoi(n.Data[1:])
				h.skipNewline = true
				parse(strings.Repeat("#", j)+" ", true)
				buf.WriteString("\n")
			case "code":
				h.skipNewline = false
				lang := h.Attr("class", n)
				var newline = ""
				if lang == "" && n.Parent != nil && n.Parent.Data == "pre" {
					lang = h.Attr("class", n.Parent)
					newline = "\n"
				}
				lang = strings.ReplaceAll(lang, "hljs ", "")
				lang = strings.ReplaceAll(lang, "highlight ", "")
				lang = strings.ReplaceAll(lang, "highlight-source-", "")
				lang = strings.ReplaceAll(lang, "language-", "")
				if lang != "" {
					newline = "\n"
				}
				buf.WriteString(newline)
				buf.WriteString("```")
				buf.WriteString(lang)
				buf.WriteString(newline)
				parse("", true)
				buf.WriteString(newline)
				buf.WriteString("```")
			case "ul", "ol":
				h.ulN++
				parse("", true)
				h.ulN--
			case "li":
				h.skipNewline = true
				buf.WriteString("\n")
				if h.ulN > 0 {
					buf.WriteString(strings.Repeat("	", h.ulN))
				}
				parse("- ", true)
			case "blockquote":
				h.skipNewline = true
				h.blockquoteN++
				buf.WriteString("\n")
				parse(strings.Repeat(">", h.blockquoteN)+" ", true)
				h.blockquoteN--
				h.skipNewline = false
			case "tr":
				if h.tdN > 0 && !h.tableSpliced {
					buf.WriteString("\n| ")
					buf.WriteString(strings.Repeat("---- | ", h.tdN))
					h.tdN = 0
					h.tableSpliced = true
				}
				buf.WriteString("\n| ")
			case "td", "th":
				parse("", true)
				buf.WriteString(" | ")
				h.tdN++
			case "pre":
				if n.FirstChild != nil && n.FirstChild.Data != "code" {
					parse("\n```\n", false)
				}
			case "p":
				if !h.skipNewline {
					buf.WriteString("\n")
				}
			}
		}
		if n != nil && n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}

	f(h.Node)

	return buf.String()
}
