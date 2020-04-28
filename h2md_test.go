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
	var table = `<table>
<thead>
<tr>
<th>CODE</th>
<th>HTTP Operation</th>
<th>Body Contents</th>
<th>Decription</th>
</tr>
</thead>
<tbody>
<tr>
<td>200</td>
<td>GET,PUT</td>
<td>资源</td>
<td>操作成功</td>
</tr>
<tr>
<td>201</td>
<td>POST</td>
<td>资源,元数据</td>
<td>对象创建成功</td>
</tr>
<tr>
<td>202</td>
<td>POST,PUT,DELETE,PATCH</td>
<td>N/A</td>
<td>请求已被接受</td>
</tr>
<tr>
<td>204</td>
<td>DELETE,PUT,PATCH</td>
<td>N/A</td>
<td>操作已经执行成功，但是没有返回结果</td>
</tr>
<tr>
<td>301</td>
<td>GET</td>
<td>link</td>
<td>资源已被移除</td>
</tr>
<tr>
<td>303</td>
<td>GET</td>
<td>link</td>
<td>重定向</td>
</tr>
<tr>
<td>304</td>
<td>GET</td>
<td>N/A</td>
<td>资源没有被修改</td>
</tr>
<tr>
<td>400</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>参数列表错误(缺少，格式不匹配)</td>
</tr>
<tr>
<td>401</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>未授权</td>
</tr>
<tr>
<td>403</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>访问受限，授权过期</td>
</tr>
<tr>
<td>404</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>资源，服务未找到</td>
</tr>
<tr>
<td>405</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>不允许的HTTP方法</td>
</tr>
<tr>
<td>409</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>资源冲突，或资源被锁定</td>
</tr>
<tr>
<td>415</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>不支持的数据(媒体)类型</td>
</tr>
<tr>
<td>429</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>请求过多被限制</td>
</tr>
<tr>
<td>500</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>系统内部错误</td>
</tr>
<tr>
<td>501</td>
<td>GET,POST,PUT,DELETE,PATCH</td>
<td>错误提示(消息)</td>
<td>接口未实现</td>
</tr>
</tbody>
</table>`
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
