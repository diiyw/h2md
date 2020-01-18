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
		{"<hr>", "---\n"},
		{"<code>code</code>", "```code```"},
		{"<pre class=\"hljs javascript\"><code>code</code></pre>", "```javascript\ncode\n```\n"},
		{"<blockquote>blockquote</blockquote>", "> blockquote\n"},
		{"<blockquote>blockquote<blockquote>sub blockquote</blockquote></blockquote>", "> blockquote\n\n>> sub blockquote\n"},

		{"<a href=\"xxx.com\">link</a>", "[link](xxx.com)"},
		{"<img src=\"xxx.jpg\" alt=\"image\"/>", "![image](xxx.jpg)"},

		{"<table><tr><th>table header</th></tr></table>", "\n| table header | "},
		{"<table><tr><th>table header</th><th>table header 1</th></tr></table>", "\n| table header | table header 1 | "},
		{
			"<table><tr><th>table header</th><th>table header 1</th></tr><tr><td>table data</td><td>table data 1</td></tr></table>",
			"\n| table header | table header 1 | \n| ---- | ---- | \n| table data| table data 1 | ",
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

func TestParseTable(t *testing.T) {
	var table = `<table>
<thead>
<tr>
<th>数据类型</th>
<th>字节长度</th>
<th>说明</th>
</tr>
</thead>
<tbody>
<tr>
<td>BOOLEAN</td>
<td>1</td>
<td>布尔值</td>
</tr>
<tr>
<td>INT8</td>
<td>1</td>
<td>单字节整型，-2^7 ~ 2^7-1</td>
</tr>
<tr>
<td>INT16</td>
<td>2</td>
<td>双字节整型，大端序，范围 -2^15 ~ 2^15 - 1</td>
</tr>
<tr>
<td>INT32</td>
<td>4</td>
<td>四字节整型、大端序，范围 -2^31 ~ 2^31 - 1</td>
</tr>
<tr>
<td>INT64</td>
<td>8</td>
<td>八字节整型、大端序，范围 -2^63 ~ 2^63 -1</td>
</tr>
<tr>
<td>UINT32</td>
<td>4</td>
<td>十字街</td>
</tr>
<tr>
<td>UUID</td>
<td>16</td>
<td>16字节，Java UUID类型</td>
</tr>
<tr>
<td>STRING</td>
<td>2+N</td>
<td>头部由2字节标识字符串长度N，后续N字节为字符串内容</td>
</tr>
<tr>
<td>NULLABLE_STRING</td>
<td>2+N</td>
<td>头部由2字节标识字符串长度N，后续N字节为字符串内容，N为-1时无后续内容</td>
</tr>
<tr>
<td>BYTES</td>
<td>4+N</td>
<td>头部4字节标识字节数组长度，后续N字节为字节数组内容</td>
</tr>
<tr>
<td>NULLABLE_BYTES</td>
<td>4+N</td>
<td>头部4字节标识字节数组长度，后续N字节为字节数组内容，N为-1时无后续内容</td>
</tr>
<tr>
<td>ARRAY</td>
<td>4+N*M</td>
<td>头部4字节标识数组长度N，M为单个数组元素的长度，N为-1时为空数组</td>
</tr>
</tbody>
</table>`
	h, err := NewH2MD(table)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(h.Text())
}
