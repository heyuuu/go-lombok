# go-lombok

## 介绍

Go 版本的属性 getter/setter 生成工具，类似 java 中的 lombok。

## 安装

推荐直接使用 go 源码安装

```bash
go install github.com/heyuuu/go-lombok@latest
```

## 常用命令

- `go-lombok generte -d {src-dir}`: 在 `src-dir` 目录(默认为当前目录)生成 getter/setter 代码，
- `go-lombok clear -d {src-dir}`: 在 `src-dir` 目录(默认为当前目录)清理生成 getter/setter 的代码

其他命令细节可通过 `go-lombok --help` 查看