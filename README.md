# H2MD

A tool to help you translate html 2 markdown.
 
几乎完美的将html转换成markdown格式的小工具，不使用正则、支持表格、代码块的转换

## Require

- go > 1.11

```go
h2md,err := NewH2md("<h1>Title</h1>")
if err == nil {
    fmt.Println(h2md.Text())
}
```

## Support tags

- a
- b
- i
- hr
- strong
- del
- img
- pre > code
- code 
- h1,h2,h3,h4,h5,h6
- table
- li
- li > ul > li
- blockquote
- blockquote > blockquote

