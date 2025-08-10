package lombok

import (
	"github.com/heyuuu/go-lombok/internal/utils/astkit"
	"go/ast"
	"go/token"
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
		if isValidIdent(prop.Getter) {
			if prop.IsRefGetter {
				getter := &ast.FuncDecl{
					Recv: recv,
					Name: ast.NewIdent(prop.Getter),
					Type: &ast.FuncType{
						Params: astkit.Fields(),
						Results: astkit.Fields(&ast.Field{
							Type: &ast.StarExpr{X: resolveTyp},
						}),
					},
					Body: astkit.BlockStmt(
						astkit.ReturnStmt(&ast.UnaryExpr{Op: token.AND, X: propFetch}),
					),
				}
				result = append(result, getter)
			} else {
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
		}

		// setter
		if isValidIdent(prop.Setter) {
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
