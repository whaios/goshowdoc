package parser

import "go/ast"

// Package 包中包含的文件和类型
type Package struct {
	Name            string // 包名，如：book
	Files           map[string]*ast.File
	TypeDefinitions map[string]*TypeSpecDef
}

// TypeSpecDef ast.TypeSpec 信息
type TypeSpecDef struct {
	PkgPath  string // 完整包名
	File     *ast.File
	TypeSpec *ast.TypeSpec
}

// Name 类型名称.
// 如：Book
func (t *TypeSpecDef) Name() string {
	if t.TypeSpec != nil {
		return t.TypeSpec.Name.Name
	}

	return ""
}

// FullName 类型全名.
// 如：book.Book
func (t *TypeSpecDef) FullName() string {
	return getFullTypeName(t.File.Name.Name, t.TypeSpec.Name.Name)
}

// FullPath 完整包名.类型名
// 如：ginweb.book.Book
func (t *TypeSpecDef) FullPath() string {
	return t.PkgPath + "." + t.Name()
}
