package lombok

import (
	f "github.com/heyuuu/go-lombok/internal/utils/astkit"
	"github.com/heyuuu/go-lombok/internal/utils/strkit"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

func ScanPkgInfo(dirPkg string, filePaths []string) (*PkgInfo, error) {
	sc := newScanner(dirPkg)
	for _, filePath := range filePaths {
		err := sc.scanFile(filePath)
		if err != nil {
			return nil, err
		}
	}
	return sc.pkg, nil
}

type scanner struct {
	pkg     *PkgInfo
	imports map[string]string
}

func newScanner(pkg string) *scanner {
	return &scanner{
		pkg: NewPkgInfo(pkg),
	}
}

func (sc *scanner) scanFile(file string) error {
	astFile, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	sc.pkg.Name = astFile.Name.Name
	sc.imports = map[string]string{}
	for _, importSpec := range astFile.Imports {
		path := strings.Trim(importSpec.Path.Value, `"`)
		var name string
		if importSpec.Name != nil {
			name = importSpec.Name.Name
		} else if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
			name = path[idx+1:]
		} else {
			name = path
		}
		sc.imports[name] = path
	}
	ast.Inspect(astFile, sc.inspectNode)
	return nil
}

func (sc *scanner) inspectNode(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.TypeSpec:
		sc.inspectTypeSpec(x)
	case *ast.FuncDecl:
		sc.inspectFuncDecl(x)
	}
	return true
}

func (sc *scanner) updateRecvName(typ string, name string) {
	if name == "" || typ == "" {
		return
	}
	sc.pkg.FindOrInitType(typ).AddExistRecvName(name)
}

func (sc *scanner) inspectFuncDecl(funcDecl *ast.FuncDecl) {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) != 1 {
		return
	}
	field := funcDecl.Recv.List[0]
	if len(field.Names) != 1 || field.Type == nil {
		return
	}

	// 检查类型使用过的 recv 名
	recvName := field.Names[0].Name
	recvType := field.Type
	if t, ok := recvType.(*ast.StarExpr); ok {
		recvType = t.X
	}
	ident, ok := recvType.(*ast.Ident)
	if !ok {
		return
	}
	recvTypeName := ident.Name
	if recvName == "" || recvName == "_" || recvTypeName == "" {
		return
	}
	typ := sc.pkg.FindOrInitType(recvTypeName)
	typ.AddExistRecvName(recvName)

	// 判断是否是事实上的 Getter/Setter
	// check getter or setter
	if funcDecl.Body == nil || len(funcDecl.Body.List) != 1 {
		return
	}
	fnType := funcDecl.Type
	stmt := funcDecl.Body.List[0]
	if len(fnType.Params.List) == 0 && fnType.Results != nil && len(fnType.Results.List) == 1 { // check getter
		retStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(retStmt.Results) != 1 {
			return
		}
		sel, ok := retStmt.Results[0].(*ast.SelectorExpr)
		if !ok {
			return
		}
		obj, ok := sel.X.(*ast.Ident)
		if !ok {
			return
		}
		if obj.Name != recvName {
			return
		}

		propName := sel.Sel.Name

		typ.FindOrInitProperty(propName).AddExistGetter(funcDecl.Name.Name)
	} else if len(fnType.Params.List) == 1 && (fnType.Results == nil || len(fnType.Results.List) == 0) {
		// check param
		if len(fnType.Params.List) != 1 {
			return
		}
		valueName := fnType.Params.List[0].Names[0].Name

		// check stmt
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 || assign.Tok != token.ASSIGN {
			return
		}
		right, ok := assign.Rhs[0].(*ast.Ident)
		if !ok || right.Name != valueName {
			return
		}
		left, ok := assign.Lhs[0].(*ast.SelectorExpr)
		if !ok {
			return
		}
		obj, ok := left.X.(*ast.Ident)
		if !ok || obj.Name != recvName {
			return
		}
		propName := left.Sel.Name

		typ.FindOrInitProperty(propName).AddExistSetter(funcDecl.Name.Name)
	}
}

func (sc *scanner) inspectTypeSpec(typeSpec *ast.TypeSpec) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	typeName := typeSpec.Name.Name
	typ := sc.pkg.FindOrInitType(typeName)

	var propertyNames []string
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			// 记录属性顺序
			propertyNames = append(propertyNames, name.Name)

			// 解析属性特性
			prop := typ.FindOrInitProperty(name.Name)
			prop.Type = sc.resolveType(field.Type)
			if field.Tag != nil {
				sc.parsePropertyTag(typ, prop, field.Tag.Value)
			}
		}
	}
	typ.PropertyNames = propertyNames
}

func (sc *scanner) parsePropertyTag(typ *Type, prop *Property, tagStr string) {
	prop.Tag = tagStr

	if len(tagStr) <= 2 {
		return
	}
	tagStr = tagStr[1 : len(tagStr)-1] // 移除头尾的 "`"

	// getter/setter from tag
	tag := reflect.StructTag(strings.Trim(tagStr, "`"))
	if recvVal, ok := tag.Lookup("recv"); ok {
		typ.RecvName = recvVal
	}
	if propVal, ok := tag.Lookup("prop"); ok {
		if propVal == "@" {
			prop.Getter = "Get" + sc.getterName(prop.Name)
			prop.Setter = sc.setterName(prop.Name)
		} else {
			if propVal == "" {
				propVal = prop.Name
			}
			prop.Getter = sc.getterName(propVal)
			prop.Setter = sc.setterName(propVal)
		}
	}
	if getterVal, ok := tag.Lookup("get"); ok {
		if getterVal == "@" {
			prop.Getter = "Get" + sc.getterName(prop.Name)
		} else if getterVal == "" {
			prop.Getter = sc.getterName(prop.Name)
		} else {
			prop.Getter = getterVal
		}
	}
	if setterVal, ok := tag.Lookup("set"); ok {
		if setterVal == "" {
			prop.Setter = sc.setterName(prop.Name)
		} else {
			prop.Setter = setterVal
		}
	}
	if prop.Getter == "" && prop.Setter == "" {
		return
	}
}

func (sc *scanner) resolveType(typ ast.Expr) ast.Expr {
	switch x := typ.(type) {
	case *ast.SelectorExpr: // p.T
		if ident, ok := x.X.(*ast.Ident); ok {
			if realPkg, ok := sc.imports[ident.Name]; ok {
				x.X = f.Ident(realPkg)
			}
		}
	case *ast.StarExpr: // *T
		x.X = sc.resolveType(x.X)
	case *ast.IndexExpr: // T1[T2]
		x.X = sc.resolveType(x.X)
		x.Index = sc.resolveType(x.Index)
	case *ast.ArrayType:
		x.Elt = sc.resolveType(x.Elt)
	}
	return typ
}

func (sc *scanner) getterName(name string) string {
	return strkit.UpperCamelCase(name)
}
func (sc *scanner) setterName(name string) string {
	return "Set" + strkit.UpperCamelCase(name)
}
