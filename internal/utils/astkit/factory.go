package astkit

import (
	"go/ast"
	"go/token"
	"strconv"
)

var (
	nilIdent   = Ident("nil")
	falseIdent = Ident("false")
	trueIdent  = Ident("true")
)

func Ident(name string) *ast.Ident { return &ast.Ident{Name: name} }

// lit
func Nil() *ast.Ident { return nilIdent }
func Int(val int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(val)}
}
func True() ast.Expr  { return trueIdent }
func False() ast.Expr { return falseIdent }
func Bool(val bool) ast.Expr {
	if val {
		return trueIdent
	}
	return falseIdent
}
func String(val string) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(val)}
}

func Array(typ ast.Expr, elements []ast.Expr) ast.Expr {
	return &ast.CompositeLit{
		Type: ArrayType(typ),
		Elts: elements,
	}
}

func Struct(fields ...ast.Expr) ast.Expr {
	return &ast.CompositeLit{Elts: fields}
}

func KeyValue(key string, value ast.Expr) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{Key: Ident(key), Value: value}
}

func Field(name *ast.Ident, typ ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{name},
		Type:  typ,
	}
}
func Fields(fields ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: fields}
}

func Type(name string) ast.Expr       { return Ident(name) }
func RefType(typ ast.Expr) ast.Expr   { return &ast.StarExpr{X: typ} }
func ArrayType(typ ast.Expr) ast.Expr { return &ast.ArrayType{Elt: typ} }

func Not(expr ast.Expr) ast.Expr { return &ast.UnaryExpr{Op: token.NOT, X: expr} }

func Call(name ast.Expr, args []ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  name,
		Args: args,
	}
}

func MethodCallExpr(instance ast.Expr, method string, args []ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   instance,
			Sel: Ident(method),
		},
		Args: args,
	}
}

func BinaryExpr(Op token.Token, first ast.Expr, others ...ast.Expr) ast.Expr {
	var expr ast.Expr = first
	for _, other := range others {
		expr = &ast.BinaryExpr{X: expr, Op: Op, Y: other}
	}
	return expr
}

func ExprStmt(expr ast.Expr) ast.Stmt {
	return &ast.ExprStmt{X: expr}
}
func AssignStmt(variable ast.Expr, value ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{variable},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{value},
	}
}
func MultiAssignStmt(variables []ast.Expr, value ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: variables,
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{value},
	}
}
func DefineStmt(variable ast.Expr, value ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{variable},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{value},
	}
}
func MultiDefineStmt(variables []ast.Expr, value ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: variables,
		Tok: token.DEFINE,
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

func VarDecl(doc *ast.CommentGroup, variable *ast.Ident, value ast.Expr) *ast.GenDecl {
	return &ast.GenDecl{
		Doc: doc,
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names:  []*ast.Ident{variable},
				Values: []ast.Expr{value},
			},
		},
	}
}
