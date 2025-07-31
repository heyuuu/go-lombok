package lombok

import (
	"github.com/heyuuu/go-lombok/internal/utils/astkit"
	"go/ast"
	"slices"
)

func GenFileCode(pkg *PkgInfo) (string, bool) {
	gen := &propertiesFileBuilder{}
	astFile := gen.generate(pkg)
	if !gen.Written() {
		return "", false
	}
	return astkit.PrintNode(astFile), true
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
	if recvName == "" && len(typ.ExistRecvNames) == 1 {
		for name, _ := range typ.ExistRecvNames {
			return name
		}
	}
	//for _, c := range []byte(typ.Name) {
	//	if ascii.IsAlpha(c) {
	//		c = ascii.ToLower(c)
	//		return string(c)
	//	}
	//}
	return "t"
}

func (b *propertiesFileBuilder) buildTypeProperties(typ *Type) []ast.Decl {
	// 过滤出需处理的属性
	properties := slices.DeleteFunc(typ.Properties(), func(property *Property) bool {
		return property.Getter == "" && property.Setter == ""
	})
	if len(properties) == 0 {
		return nil
	}

	recvName := b.getRecvName(typ)

	recv := astkit.Fields(
		astkit.Field(astkit.Ident(recvName), astkit.RefType(astkit.Ident(typ.Name))),
	)

	var result []ast.Decl

	for _, prop := range properties {
		propFetch := &ast.SelectorExpr{X: astkit.Ident(recvName), Sel: astkit.Ident(prop.Name)}
		valueName := "v"
		if valueName == recvName {
			valueName = "value"
		}
		resolveTyp := b.resolveType(prop.Type)

		if prop.Getter != "" {
			getter := &ast.FuncDecl{
				Recv: recv,
				Name: astkit.Ident(prop.Getter),
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
		if prop.Setter != "" {
			setter := &ast.FuncDecl{
				Recv: recv,
				Name: astkit.Ident(prop.Setter),
				Type: &ast.FuncType{
					Params: astkit.Fields(
						astkit.Field(astkit.Ident(valueName), resolveTyp),
					),
				},
				Body: astkit.BlockStmt(
					astkit.AssignStmt(
						propFetch,
						astkit.Ident(valueName),
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
