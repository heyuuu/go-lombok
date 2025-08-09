package lombok

import (
	"github.com/heyuuu/go-lombok/internal/utils/astkit"
	"go/ast"
)

func GenFileCode(pkg *PkgInfo) string {
	builder := &propertiesFileBuilder{}
	astFile := builder.generate(pkg)
	if !builder.Written() {
		return ""
	}
	return astkit.PrintNode(astFile)
}

type propertiesFileBuilder struct {
	*astkit.FileBuilder
}

func (b *propertiesFileBuilder) generate(pkg *PkgInfo) *ast.File {
	b.FileBuilder = astkit.NewFileBuilder(pkg.Name, pkg.Pkg)

	for _, typ := range pkg.SortedTypes() {
		for _, decl := range b.buildTypeProperties(typ) {
			b.FileBuilder.AddDecl(decl)
		}
	}

	return b.BuildFile()
}

func (b *propertiesFileBuilder) getRecvName(typ *Type) string {
	recvName := typ.RecvName
	if recvName == "" {
		if existsRecvName, ok := typ.ExistsRecvName(); ok {
			return existsRecvName
		}
	}

	return "t"
}

func (b *propertiesFileBuilder) buildTypeProperties(typ *Type) []ast.Decl {
	// build recv
	recvName := b.getRecvName(typ)
	recv := astkit.Fields(
		astkit.Field(ast.NewIdent(recvName), astkit.RefType(ast.NewIdent(typ.Name))),
	)

	var result []ast.Decl

	for prop := range typ.Properties() {
		// 跳过无需处理的属性
		if prop.Getter == "" && prop.Setter == "" {
			continue
		}

		propFetch := &ast.SelectorExpr{X: ast.NewIdent(recvName), Sel: ast.NewIdent(prop.Name)}
		valueName := "v"
		if valueName == recvName {
			valueName = "value"
		}
		resolveTyp := b.resolveType(prop.Type)

		// getter
		if b.isValidMethodName(prop.Getter) {
			getter := &ast.FuncDecl{
				Recv: recv,
				Name: ast.NewIdent(prop.Getter),
				Type: &ast.FuncType{
					Params: astkit.Fields(),
					Results: astkit.Fields(&ast.Field{
						Type: resolveTyp,
					}),
				},
				Body: astkit.BlockStmt(
					astkit.ReturnStmt(propFetch),
				),
			}
			result = append(result, getter)
		}

		// setter
		if b.isValidMethodName(prop.Setter) {
			setter := &ast.FuncDecl{
				Recv: recv,
				Name: ast.NewIdent(prop.Setter),
				Type: &ast.FuncType{
					Params: astkit.Fields(
						astkit.Field(ast.NewIdent(valueName), resolveTyp),
					),
				},
				Body: astkit.BlockStmt(
					astkit.AssignStmt(
						propFetch,
						ast.NewIdent(valueName),
					),
				),
			}
			result = append(result, setter)
		}
	}

	// 首行注释
	if len(result) > 0 {
		result[0].(*ast.FuncDecl).Doc = astkit.DocComment("\n// properties for " + typ.Name)
	}

	return result
}

func (b *propertiesFileBuilder) resolveType(typ ast.Expr) ast.Expr {
	switch x := typ.(type) {
	case *ast.SelectorExpr:
		if ident, ok := x.X.(*ast.Ident); ok {
			typ = b.PkgIdent(ident.Name, x.Sel.Name)
		}
	case *ast.StarExpr:
		x.X = b.resolveType(x.X)
	case *ast.IndexExpr:
		x.X = b.resolveType(x.X)
		x.Index = b.resolveType(x.Index)
	case *ast.ArrayType:
		x.Elt = b.resolveType(x.Elt)
	}
	return typ
}

// 判断是否是合法方法名
// 规则: [a-zA-Z_][a-zA-Z0-9_]*
func (b *propertiesFileBuilder) isValidMethodName(name string) bool {
	if name == "" {
		return false
	}

	for index, c := range []byte(name) {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_' {
			continue
		}
		if index > 0 && ('0' <= c && c <= '9') {
			continue
		}
		return false
	}

	return true
}
