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
	replacer map[string]func(val string, n *html.Node) string
}

// NewH2MD create H2MD with html text
func NewH2MD(htmlText string) (*H2MD, error) {
	node, err := html.Parse(strings.NewReader(htmlText))
	if err == nil {
		return &H2MD{node, make(map[string]func(val string, n *html.Node) string)}, nil
	}
	return nil, err
}

//NewH2MDFromNode create H2MD with html node
func NewH2MDFromNode(node *html.Node) (*H2MD, error) {
	return &H2MD{node, make(map[string]func(val string, n *html.Node) string)}, nil
}

// Replace Replace element attribute value
func (h *H2MD) Replace(attr string, r func(val string, n *html.Node) string) {
	h.replacer[attr] = r
}

// Attr Return the element attribute
func (h *H2MD) Attr(name string, n *html.Node) string {
	for _, attr := range n.Attr {
		if name == attr.Key {
			if r, ok := h.replacer[name]; ok {
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

	var tableColumn int

	var spitedTable bool

	var tdCounter int

	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			var data string
			if n.Parent != nil {
				switch n.Parent.Data {
				case "a":
					data = "[" + n.Data + "](" + h.Attr("href", n.Parent) + ")"
				case "strong", "b":
					data = "**" + n.Data + "**"
				case "del":
					data = "~~" + n.Data + "~~"
				case "h1", "h2", "h3", "h4", "h5", "h6":
					j, _ := strconv.Atoi(n.Parent.Data[1:])
					title := ""
					for i := 0; i < j; i++ {
						title += "#"
					}
					data += title + " " + n.Data
				case "blockquote":
					if n.PrevSibling == nil {
						data += "> "
					}
					data += n.Data + "\n"
				case "code":
					lang := h.Attr("class", n.Parent)
					var newline string
					if lang == "" && n.Parent.Parent != nil && n.Parent.Parent.Data == "pre" {
						class := h.Attr("class", n.Parent.Parent)
						newline = "\n"
						lang = strings.ReplaceAll(class, "hljs ", "")
						lang = strings.ReplaceAll(lang, "highlight ", "")
						lang = strings.ReplaceAll(lang, "highlight-source-", "")
					}
					if lang != "" {
						newline = "\n"
					}
					data = "```" + lang + newline + n.Data + newline + "```" + newline
				case "li":
					if n.PrevSibling == nil {
						data += "- "
					}
					data += n.Data
				case "i":
					data = "*" + n.Data + "*"
				case "th":
					tableColumn++
					data = n.Data + " | "
				case "td":
					if !spitedTable {
						buf.WriteString("\n| ")
						for i := 0; i < tableColumn; i++ {
							buf.WriteString("---- | ")
						}
						buf.WriteString("\n| ")
						spitedTable = true
					}
					tdCounter++
					if tdCounter == tableColumn {
						data = n.Data + " | "
						tdCounter = 0
						break
					}
					data = n.Data
				default:
					data = skipInvalidString(n)
				}
			}
			buf.WriteString(data)
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "hr":
				buf.WriteString("---\n")
			case "img":
				data := "![" + h.Attr("alt", n) + "](" + h.Attr("src", n) + ")"
				buf.WriteString(data)
			case "ul":
				if n.Parent != nil && n.Parent.Data == "li" {
					buf.WriteString("\n	")
				}
			case "blockquote":
				if n.Parent != nil && n.Parent.Data == "blockquote" {
					buf.WriteString("\n>")
				}
			case "table":
				buf.WriteString("\n| ")
			case "td":
				if spitedTable {
					buf.WriteString("| ")
				}
			case "p","li":
				buf.WriteString("\n")
			}
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}

	f(h.Node)

	return buf.String()
}

func skipInvalidString(n *html.Node) string {
	var trimNewline = []string{
		"table",
		"tr",
		"th",
		"thead",
	}
	for _, el := range trimNewline {
		if n.Parent != nil && n.Parent.Data == el {
			return ""
		}
	}
	return n.Data
}
