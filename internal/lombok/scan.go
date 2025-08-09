package lombok

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// ScanPkgInfo 扫描一个包下的所有文件，返回包信息
// 隐式要求 srcFiles 必须在同一包下
func ScanPkgInfo(pkg string, srcFiles []string) (*PkgInfo, error) {
	sc := newScanner(pkg)
	for _, srcFile := range srcFiles {
		err := sc.scanFile(srcFile)
		if err != nil {
			return nil, err
		}
	}
	return sc.pkg, nil
}

func ScanCode(pkg string, code string) (*PkgInfo, error) {
	sc := newScanner(pkg)
	err := sc.scanFileCode(code)
	if err != nil {
		return nil, err
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

	return sc.scanAstFile(astFile)
}

func (sc *scanner) scanFileCode(code string) error {
	astFile, err := parser.ParseFile(token.NewFileSet(), "", code, parser.ParseComments)
	if err != nil {
		return err
	}

	return sc.scanAstFile(astFile)
}

func (sc *scanner) scanAstFile(astFile *ast.File) error {
	sc.pkg.Name = astFile.Name.Name

	// 记录当前文件的 imports 表
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

	// 遍历文件节点
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

// 分析 struct 类型定义获取属性信息

func (sc *scanner) inspectTypeSpec(typeSpec *ast.TypeSpec) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	typeName := typeSpec.Name.Name
	typ := sc.pkg.FindOrInitType(typeName)

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			prop := typ.AddProperty(name.Name)
			prop.Type = sc.resolveType(field.Type)
			if field.Tag != nil {
				sc.parsePropertyTag(typ, prop, field.Tag.Value)
			}
		}
	}
}

func (sc *scanner) parsePropertyTag(typ *Type, prop *Property, tagStr string) {
	prop.Tag = tagStr

	if len(tagStr) <= 2 {
		return
	}

	// getter/setter from tag
	tag := reflect.StructTag(strings.Trim(tagStr, "`"))
	if recvVal, ok := tag.Lookup("recv"); ok {
		typ.RecvName = recvVal
	}
	if propVal, ok := tag.Lookup("prop"); ok {
		prop.Getter = sc.calcGetterName(prop, propVal)
		prop.Setter = sc.calcSetterName(prop, propVal, true)
	}
	if getterVal, ok := tag.Lookup("get"); ok {
		prop.Getter = sc.calcGetterName(prop, getterVal)
	}
	if setterVal, ok := tag.Lookup("set"); ok {
		prop.Setter = sc.calcSetterName(prop, setterVal, false)
	}
}

func (sc *scanner) resolveType(typ ast.Expr) ast.Expr {
	switch x := typ.(type) {
	case *ast.SelectorExpr: // p.T
		if ident, ok := x.X.(*ast.Ident); ok {
			if realPkg, ok := sc.imports[ident.Name]; ok {
				x.X = ast.NewIdent(realPkg)
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

func (sc *scanner) calcGetterName(prop *Property, tagVal string) string {
	if tagVal == "" {
		return pascalCase(prop.Name)
	} else if tagVal == "@" {
		return "Get" + pascalCase(prop.Name)
	} else {
		return tagVal
	}
}

func (sc *scanner) calcSetterName(prop *Property, tagVal string, isPropTag bool) string {
	if tagVal == "" || tagVal == "@" {
		return "Set" + pascalCase(prop.Name)
	} else if isPropTag {
		return "Set" + tagVal
	} else {
		return tagVal
	}
}

// 分析函数定义判断是否为某属性的 getter/setter

func (sc *scanner) inspectFuncDecl(funcDecl *ast.FuncDecl) {
	// 获取并检查 recv
	recvName, recvTypeName, ok := sc.getRecvOfFunc(funcDecl)
	if !ok || recvName == "" || recvName == "_" || recvTypeName == "" {
		return
	}

	// 记录使用到的 recvName
	sc.recordTypeRecvName(recvTypeName, recvName)

	// 判断是否为事实上的 getter/setter
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

		sc.recordGetter(recvTypeName, propName, funcDecl.Name.Name)
	} else if len(fnType.Params.List) == 1 && (fnType.Results == nil || len(fnType.Results.List) == 0) { // check setter
		// check param
		if len(fnType.Params.List[0].Names) != 1 {
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

		sc.recordSetter(recvTypeName, propName, funcDecl.Name.Name)
	}
}

func (sc *scanner) getRecvOfFunc(funcDecl *ast.FuncDecl) (recvName string, recvTypeName string, ok bool) {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) != 1 {
		return
	}

	field := funcDecl.Recv.List[0]
	if len(field.Names) != 1 || field.Type == nil {
		return
	}

	// 获取并检查 recv 名
	recvName = field.Names[0].Name

	// 获取并检查 recv 类型名
	recvType := field.Type
	if t, ok := recvType.(*ast.StarExpr); ok {
		recvType = t.X
	}
	ident, ok := recvType.(*ast.Ident)
	if !ok {
		return
	}
	recvTypeName = ident.Name

	return recvName, recvTypeName, true
}

func (sc *scanner) recordTypeRecvName(typName string, recvName string) {
	typ := sc.pkg.FindOrInitType(typName)
	typ.RecordExistsRecvName(recvName)
}

func (sc *scanner) recordGetter(typName string, propName string, getter string) {
	typ := sc.pkg.FindOrInitType(typName)
	prop := typ.FindOrInitProperty(propName)
	prop.RecordExistingGetter(getter)
}

func (sc *scanner) recordSetter(typName string, propName string, setter string) {
	typ := sc.pkg.FindOrInitType(typName)
	prop := typ.FindOrInitProperty(propName)
	prop.RecordExistingSetter(setter)
}
