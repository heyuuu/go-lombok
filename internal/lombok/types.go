package lombok

import (
	"github.com/heyuuu/go-lombok/internal/utils/mapkit"
	"go/ast"
	"strings"
)

type PkgInfo struct {
	Name  string
	Pkg   string
	Types map[string]*Type
}

func NewPkgInfo(pkg string) *PkgInfo {
	name := pkg
	if idx := strings.LastIndexByte(name, '/'); idx >= 0 {
		name = name[idx+1:]
	}
	return &PkgInfo{
		Name:  name,
		Pkg:   pkg,
		Types: make(map[string]*Type),
	}
}

func (pkg *PkgInfo) FindOrInitType(typeName string) *Type {
	if t, exists := pkg.Types[typeName]; exists {
		return t
	} else {
		typ := NewType(typeName)
		pkg.Types[typeName] = typ
		return typ
	}
}

func (pkg *PkgInfo) SortedTypes() []*Type {
	return mapkit.SortedValues(pkg.Types)
}

type Type struct {
	Name           string
	RecvName       string
	PropertyNames  []string // 属性名列表，按类型定义字段顺序
	PropertyMap    map[string]*Property
	ExistRecvNames map[string]bool
}

func NewType(name string) *Type {
	return &Type{
		Name:           name,
		PropertyMap:    make(map[string]*Property),
		ExistRecvNames: make(map[string]bool),
	}
}

func (typ *Type) FindOrInitProperty(propName string) *Property {
	if p, exists := typ.PropertyMap[propName]; exists {
		return p
	} else {
		prop := NewProperty(propName)
		typ.PropertyMap[propName] = prop
		return prop
	}
}

func (typ *Type) Properties() []*Property {
	properties := make([]*Property, len(typ.PropertyNames))
	for i, name := range typ.PropertyNames {
		properties[i] = typ.FindOrInitProperty(name)
	}
	return properties
}

func (typ *Type) AddExistRecvName(name string) {
	typ.ExistRecvNames[name] = true
}

type Property struct {
	Name         string
	Getter       string
	Setter       string
	Tag          string
	Type         ast.Expr
	ExistGetters map[string]bool
	ExistSetters map[string]bool
}

func NewProperty(name string) *Property {
	return &Property{
		Name:         name,
		ExistGetters: make(map[string]bool),
		ExistSetters: make(map[string]bool),
	}
}

func (typ *Property) AddExistGetter(name string) {
	typ.ExistGetters[name] = true
}

func (typ *Property) AddExistSetter(name string) {
	typ.ExistSetters[name] = true
}
