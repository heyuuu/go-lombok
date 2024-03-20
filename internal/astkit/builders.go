package astkit

import (
	"github.com/heyuuu/go-lombok/internal/mapkit"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// Imports
type Imports struct {
	imports map[string]string
	alias   map[string]bool
}

func (imports *Imports) init() {
	if imports.imports == nil {
		imports.imports = map[string]string{}
		imports.alias = map[string]bool{}
	}
}

func (imports *Imports) FindOrAdd(pkgName string) string {
	imports.init()

	pkgName = strings.Trim(pkgName, "/")
	if alias, ok := imports.imports[pkgName]; ok {
		return alias
	} else {
		newAlias := imports.newImportAlias(pkgName)
		imports.imports[pkgName] = newAlias
		imports.alias[newAlias] = true
		return newAlias
	}
}

func (imports *Imports) Build() *ast.GenDecl {
	if len(imports.imports) == 0 {
		return nil
	}

	pkgNames := mapkit.SortedKeys(imports.imports)

	importSpecs := make([]ast.Spec, len(pkgNames))
	for i, pkgName := range pkgNames {
		aliasName := imports.imports[pkgName]

		var aliasNameNode *ast.Ident
		if aliasName != imports.getDefaultAlias(pkgName) {
			aliasNameNode = Ident(aliasName)
		}

		importSpecs[i] = &ast.ImportSpec{
			Path: String(pkgName),
			Name: aliasNameNode,
		}
	}

	return &ast.GenDecl{Tok: token.IMPORT, Specs: importSpecs}
}

func (imports *Imports) getDefaultAlias(pkgName string) string {
	if idx := strings.LastIndexByte(pkgName, '/'); idx >= 0 {
		return pkgName[idx+1:]
	}
	return pkgName
}

func (imports *Imports) newImportAlias(pkgName string) string {
	alias := imports.getDefaultAlias(pkgName)
	if !imports.alias[alias] {
		return alias
	}

	for i := 2; ; i++ {
		newAlias := alias + strconv.Itoa(i)
		if !imports.alias[newAlias] {
			return newAlias
		}
	}
}

// FileBuilder
type FileBuilder struct {
	name    string
	pkg     string
	imports Imports
	decls   []ast.Decl
}

func NewFileBuilder(name string, pkg string) *FileBuilder {
	return &FileBuilder{name: name, pkg: pkg}
}

func (b *FileBuilder) Written() bool {
	return len(b.decls) != 0
}

func (b *FileBuilder) AddDecl(decl ast.Decl) {
	b.decls = append(b.decls, decl)
}

func (b *FileBuilder) BuildFile() *ast.File {
	var decls []ast.Decl
	if importDecl := b.imports.Build(); importDecl != nil {
		decls = append(decls, importDecl)
	}
	decls = append(decls, b.decls...)

	return &ast.File{
		Name:  Ident(b.name),
		Decls: decls,
	}
}

// expr
func (b *FileBuilder) PkgIdent(pkg string, name string) ast.Expr {
	if pkg == b.pkg {
		return &ast.Ident{Name: name}
	}

	pkgAlias := b.imports.FindOrAdd(pkg)
	return &ast.SelectorExpr{
		X:   &ast.Ident{Name: pkgAlias},
		Sel: &ast.Ident{Name: name},
	}
}
