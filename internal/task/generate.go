package task

import (
	f "github.com/heyuuu/go-lombok/internal/astkit"
	"go/ast"
	"slices"
)

func GenFileCode(pkg *PkgInfo) (string, bool) {
	gen := &propertiesFileBuilder{}
	astFile := gen.generate(pkg)
	if !gen.Written() {
		return "", false
	}
	return f.PrintNode(astFile), true
}

type propertiesFileBuilder struct {
	*f.FileBuilder
}

func (b *propertiesFileBuilder) generate(pkg *PkgInfo) *ast.File {
	b.FileBuilder = f.NewFileBuilder(pkg.Name, pkg.Pkg)

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

	recv := f.Fields(
		f.Field(f.Ident(recvName), f.RefType(f.Ident(typ.Name))),
	)

	var result []ast.Decl

	for _, prop := range properties {
		propFetch := &ast.SelectorExpr{X: f.Ident(recvName), Sel: f.Ident(prop.Name)}
		valueName := "v"
		if valueName == recvName {
			valueName = "value"
		}
		resolveTyp := b.resolveType(prop.Type)

		if prop.Getter != "" {
			getter := &ast.FuncDecl{
				Recv: recv,
				Name: f.Ident(prop.Getter),
				Type: &ast.FuncType{
					Params: f.Fields(),
					Results: f.Fields(&ast.Field{
						Type: resolveTyp,
					}),
				},
				Body: f.BlockStmt(
					f.ReturnStmt(propFetch),
				),
			}
			result = append(result, getter)
		}
		if prop.Setter != "" {
			setter := &ast.FuncDecl{
				Recv: recv,
				Name: f.Ident(prop.Setter),
				Type: &ast.FuncType{
					Params: f.Fields(
						f.Field(f.Ident(valueName), resolveTyp),
					),
				},
				Body: f.BlockStmt(
					f.AssignStmt(
						propFetch,
						f.Ident(valueName),
					),
				),
			}
			result = append(result, setter)
		}
	}

	// 首行注释
	if len(result) > 0 {
		result[0].(*ast.FuncDecl).Doc = f.DocComment("\n// properties for " + typ.Name)
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
