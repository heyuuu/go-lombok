package lombok

import (
	"go/ast"
	"iter"
	"maps"
	"slices"
	"strings"
)

type PkgInfo struct {
	Name string
	Pkg  string
	// private
	types map[string]*Type
}

func NewPkgInfo(pkg string) *PkgInfo {
	name := pkg
	if idx := strings.LastIndexByte(name, '/'); idx >= 0 {
		name = name[idx+1:]
	}
	return &PkgInfo{
		Name:  name,
		Pkg:   pkg,
		types: make(map[string]*Type),
	}
}

func (pkg *PkgInfo) FindOrInitType(typeName string) *Type {
	if t, exists := pkg.types[typeName]; exists {
		return t
	} else {
		typ := NewType(typeName)
		pkg.types[typeName] = typ
		return typ
	}
}

func (pkg *PkgInfo) SortedTypes() []*Type {
	types := slices.Collect(maps.Values(pkg.types))
	slices.SortFunc(types, func(a, b *Type) int {
		return strings.Compare(a.Name, b.Name)
	})
	return types
}

type Type struct {
	Name     string
	RecvName string
	// private
	propertyNames   []string // 属性名列表，按类型定义字段顺序
	propertyMap     map[string]*Property
	existsRecvNames map[string]bool // 已存在的 recv 名
}

func NewType(name string) *Type {
	return &Type{
		Name:            name,
		propertyNames:   nil,
		propertyMap:     make(map[string]*Property),
		existsRecvNames: make(map[string]bool),
	}
}

func (typ *Type) AddProperty(propName string) *Property {
	typ.propertyNames = append(typ.propertyNames, propName)
	return typ.FindOrInitProperty(propName)
}

func (typ *Type) FindOrInitProperty(propName string) *Property {
	if p, exists := typ.propertyMap[propName]; exists {
		return p
	} else {
		prop := NewProperty(propName)
		typ.propertyMap[propName] = prop
		return prop
	}
}

func (typ *Type) Properties() iter.Seq[*Property] {
	return func(yield func(*Property) bool) {
		for _, propName := range typ.propertyNames {
			prop := typ.propertyMap[propName]
			if !yield(prop) {
				return
			}
		}
	}
}

func (typ *Type) RecordExistsRecvName(name string) {
	typ.existsRecvNames[name] = true
}

func (typ *Type) ExistsRecvName() (string, bool) {
	if len(typ.existsRecvNames) == 1 {
		for name, _ := range typ.existsRecvNames {
			return name, true
		}
	}
	return "", false
}

type Property struct {
	Name        string
	Getter      string
	IsRefGetter bool
	Setter      string
	Tag         string
	Type        ast.Expr

	// private
	existingGetters []string
	existingSetters []string
}

func NewProperty(name string) *Property {
	return &Property{
		Name: name,
	}
}

func (prop *Property) ExistsGetter(name string) bool {
	for _, getter := range prop.existingGetters {
		if getter == name {
			return true
		}
	}
	return false
}

func (prop *Property) RecordExistingGetter(name string) {
	if prop.ExistsGetter(name) {
		return
	}
	prop.existingGetters = append(prop.existingGetters, name)
}

func (prop *Property) ExistingGetters() iter.Seq[string] {
	return slices.Values(prop.existingGetters)
}

func (prop *Property) ExistsSetter(name string) bool {
	for _, setter := range prop.existingSetters {
		if setter == name {
			return true
		}
	}
	return false
}

func (prop *Property) RecordExistingSetter(name string) {
	if prop.ExistsSetter(name) {
		return
	}
	prop.existingSetters = append(prop.existingSetters, name)
}

func (prop *Property) ExistingSetters() iter.Seq[string] {
	return slices.Values(prop.existingSetters)
}
