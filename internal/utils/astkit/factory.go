package astkit

import (
	"go/ast"
	"go/token"
)

func Field(name *ast.Ident, typ ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{name},
		Type:  typ,
	}
}
func Fields(fields ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: fields}
}

func RefType(typ ast.Expr) ast.Expr { return &ast.StarExpr{X: typ} }
func AssignStmt(variable ast.Expr, value ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{variable},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{value},
	}
}
func BlockStmt(list ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: list}
}

func ReturnStmt(results ...ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{
		Results: results,
	}
}

func DocComment(comments ...string) *ast.CommentGroup {
	if len(comments) == 0 {
		return nil
	}

	var list []*ast.Comment
	for _, comment := range comments {
		list = append(list, &ast.Comment{
			Text: comment,
		})
	}
	return &ast.CommentGroup{List: list}
}
