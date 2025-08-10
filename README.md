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

## tag 语法规则

新增支持四种tag: `get` / `set` / `prop` / `recv`

### `get`

支持值有几种情况
- `""`：生成的 Getter 函数名为 `大驼峰(属性名)`
- `"@"`：生成的 Getter 函数名为 `Get + 大驼峰(属性名)`
- `"合法函数名"`：生成的 Getter 函数名为对应函数名
- 以 `&` 为前缀，后接以上任意值：生成的 Getter 函数返回的是对应属性的引用

### `set`

支持值有几种情况
- `""`：生成的 Getter 函数名为 `Set + 大驼峰(属性名)`
- `"合法函数名"`：生成的 Getter 函数名为对应函数名

### `prop`

`prop` 支持一些常见的 `get` + `set` 组合的简写，属于语法糖；**同一属性有 `prop` 标签时不可同时出现 `get` 或 `set` 标签**

支持值有几种情况:
- `""`: 等价于 `get:"" set:""`
- `"@"`: 等价于 `get:"@" set:""`
- `"合法属性名"`: 等价于 `get:"大驼峰(属性名)" set:"Set + 大驼峰(属性名)"`
- `"@合法属性名"`: 等价于 `get:"Get + 大驼峰(属性名)" set:"Set + 大驼峰(属性名)"`
- 以 `&` 为前缀，后接以上任意值：生成的 Getter 函数返回的是对应属性的引用

### `recv`

`recv` 用于指定 getter / setter 的 recv 变量名，未指定时默认 recv 名为 `t`