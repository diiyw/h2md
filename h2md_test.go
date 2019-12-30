package h2md

import (
	"testing"
)

func TestNewH2MD(t *testing.T) {
	htmlTexts := []struct {
		text   string
		expect string
	}{
		{"<h1>Title 1</h1>", "# Title 1"},
		{"<h2>Title 2</h2>", "## Title 2"},
		{"<h3>Title 3</h3>", "### Title 3"},
		{"<h4>Title 4</h4>", "#### Title 4"},
		{"<h5>Title 5</h5>", "##### Title 5"},
		{"<h6>Title 6</h6>", "###### Title 6"},

		{"<li>List</li>", "- List"},
		{"<li>List <a href=\"xxx.com\">link</a></li>", "- List [link](xxx.com)"},
		{"<li>List <strong>strong</strong></li>", "- List **strong**"},

		{"<li>List<ul><li>sub list</li><ul></li>", "- List\n	- sub list"},

		{"<b>List</b>", "**List**"},
		{"<strong>strong</strong>", "**strong**"},
		{"<i>List</i>", "*List*"},
		{"<code>code</code>", "```code```"},
		{"<pre class=\"hljs javascript\"><code>code</code></pre>", "```javascript\ncode\n```\n"},
		{"<blockquote>blockquote</blockquote>", "> blockquote"},
		{"<blockquote>blockquote<blockquote>sub blockquote</blockquote></blockquote>", "> blockquote\n>> sub blockquote"},

		{"<a href=\"xxx.com\">link</a>", "[link](xxx.com)"},
		{"<img src=\"xxx.jpg\" alt=\"image\"/>", "![image](xxx.jpg)"},

		{"<table><th>table header</th></table>", "| table header | "},
		{"<table><th>table header</th><th>table header 1</th></table>", "| table header | table header 1 | "},
		{
			"<table><tr><th>table header</th><th>table header 1</th></tr><tr><td>table data</td><td>table data 1</td></tr></table>",
			"| table header | table header 1 | \n| ---- | ---- | \n| table data | table data 1 | \n",
		},
	}

	for _, htmlText := range htmlTexts {
		h, err := NewH2MD(htmlText.text)
		if err != nil {
			t.Error(err)
		}
		if h.Text() != htmlText.expect {
			t.Errorf("Expect \"%s\" but got \"%s\"", htmlText.expect, h.Text())
		}
	}

}
