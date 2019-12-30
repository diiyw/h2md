# H2MD

html 2 markdown. 几乎完美的将html转换成markdown格式的小工具,支持表格转换

## Require

- go > 1.11

```go
h2md,err := NewH2md()
if err == nil {
    fmt.Println(h2md.Text())
}
```

## Support tags

- a
- b
- i
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

