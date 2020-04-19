package h2md

import (
	"fmt"
	"testing"
)

func TestNewH2MD(t *testing.T) {
	htmlTexts := []struct {
		text   string
		expect string
	}{
		{"<h1>Title 1</h1>", "\n# Title 1\n"},
		{"<h2>Title 2</h2>", "\n## Title 2\n"},
		{"<h3>Title 3</h3>", "\n### Title 3\n"},
		{"<h4>Title 4</h4>", "\n#### Title 4\n"},
		{"<h5>Title 5</h5>", "\n##### Title 5\n"},
		{"<h6>Title 6</h6>", "\n###### Title 6\n"},
		{`<h1><strong>1</strong><strong>、前言</strong></h1>`, "\n# **1****、前言**\n"},
		{"<ul><li>List</li></ul>", "\n- List"},
		{"<ul><li>List <a href=\"xxx.com\">link</a></li></ul>", "\n- List [link](xxx.com)"},
		{"<ul><li>List <strong>strong</strong></li></ul>", "\n- List **strong**"},

		{"<ul><li>List<ul><li>sub list</li></ul></li></ul>", "\n- List\n	- sub list"},
		{
			"<ul><li>List<ul><li>sub list</li><li>sub2 list</li></ul></li></ul>",
			"\n- List\n	- sub list\n	- sub2 list",
		},
		{"<b>List</b>", "**List**"},
		{"<strong>strong</strong>", "**strong**"},
		{"<i>List</i>", "*List*"},
		{"<hr>", "\n---\n"},
		{"<code>code</code>", "```code```"},
		{"<pre class=\"hljs javascript\"><code>code</code></pre>", "\n```javascript\ncode\n```"},
		{"<blockquote>blockquote</blockquote>", "\n> blockquote"},
		{"<blockquote>blockquote<blockquote>sub blockquote</blockquote></blockquote>", "\n> blockquote\n>> sub blockquote"},

		{"<a href=\"xxx.com\">link</a>", "[link](xxx.com)"},
		{"<img src=\"xxx.jpg\" alt=\"image\"/>", "![image](xxx.jpg)"},

		{"<table><tr><th>table header</th></tr></table>", "\n| table header | "},
		{"<table><tr><th>table header</th><th>table header 1</th></tr></table>", "\n| table header | table header 1 | "},
		{
			"<table><tr><th>table header</th><th>table header 1</th></tr><tr><td>table data</td><td>table data 1</td></tr></table>",
			"\n| table header | table header 1 | \n| ---- | ---- | \n| table data | table data 1 | ",
		},
	}

	for _, htmlText := range htmlTexts {
		h, err := NewH2MD(htmlText.text)
		if err != nil {
			t.Error(err)
		}
		text := h.Text()
		if text != htmlText.expect {
			t.Errorf("Expect \"%s\" but got \"%s\"", htmlText.expect, text)
		}
	}
}

func TestParseTable(t *testing.T) {
	var table = `<table class="table table-bordered table-striped">
<tbody><tr>
<th>type</th>
<th>用途</th>
</tr>
<tr>
<td>Data api</td>
<td>查询文档（能被输出的规则、虚拟文档等）</td>
</tr>
<tr>
<td>Policy api</td>
<td>查询策略</td>
</tr>
<tr>
<td>Query api</td>
<td>执行命令</td>
</tr>
<tr>
<td>Compile api</td>
<td>执行部分查询计算（<code>partial evaluate query</code>）</td>
</tr>
<tr>
<td>Health api</td>
<td>健康检查</td>
</tr>
<tr>
<td>Metric api</td>
<td>指标统计（<code>prometheus</code>格式)</td>
</tr>
</tbody></table>`
	h, err := NewH2MD(table)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(h.Text())
}

func TestParseBlockquote(t *testing.T) {
	var table = `<blockquote>
<p>block1<br>
block1</p>
<blockquote>
<p>block2</p>
<blockquote>
<p>block3</p>
</blockquote>
</blockquote>
</blockquote>`
	h, err := NewH2MD(table)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(h.Text())
}
