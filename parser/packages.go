package parser

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/whaios/goshowdoc/log"
	"golang.org/x/tools/go/packages"
)

// AstFileInfo ast.File 文件信息.
type AstFileInfo struct {
	File        *ast.File
	Path        string // Go 源码文件的绝对路径
	PackagePath string // Go 源码文件完整包名
}

func NewPackages() *Packages {
	return &Packages{
		files:             make(map[*ast.File]*AstFileInfo),
		packages:          make(map[string]*Package),
		uniqueDefinitions: make(map[string]*TypeSpecDef),
	}
}

// Packages 存储扫描到的 Go 文件、包路径和他们之间的关系。
// 主要用于解析注释中可能会引用到的类型。
type Packages struct {
	projectDir string // 目标Go项目所在目录（支持绝对路径和相对路径），用于加载外部包时使用，为空时默认为当前运行目录。
	files      map[*ast.File]*AstFileInfo

	packages map[string]*Package // key=完整包名. 如：ginweb/book
	// 在目录下收集到的具有唯一名称（包名+类型名）的类型。如果存在同名的不会出现到该字典中。
	uniqueDefinitions map[string]*TypeSpecDef // key=类型全名. 如：book.Book
}

// CollectAstFile 收集 Go 源码文件
//
// @param packageDir 如：refstruct/employee
// @param absPath 如：文件绝对路径
func (p *Packages) CollectAstFile(packageDir, absPath string, astFile *ast.File) {
	log.Debug("收集Go文件: %s", absPath)
	p.files[astFile] = &AstFileInfo{
		File:        astFile,
		Path:        absPath,
		PackagePath: packageDir,
	}
}

// SortedFiles 获取所有文件，根据文件路径按字母先后顺序排序
func (p *Packages) SortedFiles() []*AstFileInfo {
	sortedFiles := make([]*AstFileInfo, 0, len(p.files))
	for _, info := range p.files {
		sortedFiles = append(sortedFiles, info)
	}

	sort.Slice(sortedFiles, func(i, j int) bool {
		return strings.Compare(sortedFiles[i].Path, sortedFiles[j].Path) < 0
	})
	return sortedFiles
}

// ParseTypes 解析所有代码文件中的类型
func (p *Packages) ParseTypes() {
	for astFile, info := range p.files {
		p.parseTypesFromFile(astFile, info.PackagePath)
	}
}

func (p *Packages) parseTypesFromFile(astFile *ast.File, pkgPath string) {
	for _, astDeclaration := range astFile.Decls {
		if generalDeclaration, ok := astDeclaration.(*ast.GenDecl); ok && generalDeclaration.Tok == token.TYPE {
			for _, astSpec := range generalDeclaration.Specs {
				if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
					typeSpecDef := &TypeSpecDef{
						PkgPath:  pkgPath,
						File:     astFile,
						TypeSpec: typeSpec,
					}

					fullName := typeSpecDef.FullName()
					log.Debug("收集类型: %s", fullName)
					anotherTypeDef, ok := p.uniqueDefinitions[fullName]
					if ok {
						if typeSpecDef.PkgPath == anotherTypeDef.PkgPath {
							continue
						} else {
							delete(p.uniqueDefinitions, fullName)
						}
					} else {
						p.uniqueDefinitions[fullName] = typeSpecDef
					}

					if p.packages[typeSpecDef.PkgPath] == nil {
						p.packages[typeSpecDef.PkgPath] = &Package{
							Name:            astFile.Name.Name,
							TypeDefinitions: map[string]*TypeSpecDef{typeSpecDef.Name(): typeSpecDef},
						}
					} else if _, ok = p.packages[typeSpecDef.PkgPath].TypeDefinitions[typeSpecDef.Name()]; !ok {
						p.packages[typeSpecDef.PkgPath].TypeDefinitions[typeSpecDef.Name()] = typeSpecDef
					}
				}
			}
		}
	}
}

// FindTypeSpec 查找类型
//
// @param fullTypeName 包名.类型名，如：ListRsp 或 book.Book
func (p *Packages) FindTypeSpec(fullTypeName string, file *ast.File) *TypeSpecDef {
	if file == nil { // for test
		return p.uniqueDefinitions[fullTypeName]
	}
	var pkgName, typeName string
	{
		typeName = fullTypeName
		parts := strings.SplitN(fullTypeName, ".", 2)
		if len(parts) == 2 {
			pkgName = parts[0]
			typeName = parts[1]
		}
	}

	// 有包名
	if pkgName != "" {
		// 从文件中导入的包中查找指定包路径
		pkgPath, isAliasPkgName := p.findPackagePathFromImports(pkgName, file)

		// 没有别名的类型，首先在唯一类型中查找。
		if !isAliasPkgName {
			if typeDef, ok := p.uniqueDefinitions[fullTypeName]; ok {
				return typeDef
			}
		}
		// 没有找到对应的包名
		if pkgPath == "" {
			return nil
		}

		// 收集外部包
		p.loadExternalPackage(pkgPath)

		return p.findTypeSpec(pkgPath, typeName)
	}

	typeDef, ok := p.uniqueDefinitions[getFullTypeName(file.Name.Name, typeName)]
	if ok {
		return typeDef
	}

	typeDef = p.findTypeSpec(p.files[file].PackagePath, typeName)
	if typeDef != nil {
		return typeDef
	}

	for _, imp := range file.Imports {
		if imp.Name != nil && imp.Name.Name == "." {
			pkgPath := strings.Trim(imp.Path.Value, `"`)
			// 收集外部包
			p.loadExternalPackage(pkgPath)

			if typeDef = p.findTypeSpec(pkgPath, typeName); typeDef != nil {
				return typeDef
			}
		}
	}

	return nil
}

// findTypeSpec 查找类型
//
// @param pkgPath 如：refstruct/employee
// @param typeName 如：Employee
func (p *Packages) findTypeSpec(pkgPath string, typeName string) *TypeSpecDef {
	if p.packages == nil {
		return nil
	}
	pd, found := p.packages[pkgPath]
	if found {
		typeSpec, ok := pd.TypeDefinitions[typeName]
		if ok {
			return typeSpec
		}
	}

	return nil
}

func (p *Packages) loadExternalPackage(importPath string) {
	if p.packages != nil {
		if _, ok := p.packages[importPath]; ok {
			// 已经收集过该包
			return
		}
	}
	log.Debug("加载外部包: %s", importPath)

	cfg := &packages.Config{
		Dir:  p.projectDir,
		Mode: packages.NeedImports | packages.NeedTypes | packages.NeedSyntax,
	}
	pkgs, _ := packages.Load(cfg, importPath)
	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			p.parseTypesFromFile(astFile, pkg.ID)
		}
	}
}

// 从文件中导入的包中查找指定包路径。
// @pkgName 包名
func (p *Packages) findPackagePathFromImports(pkgName string, file *ast.File) (pkgPath string, isAliasPkgName bool) {
	for _, imp := range file.Imports {
		// 有别名的包，别名相同，直接取得该包路径
		if imp.Name != nil && imp.Name.Name == pkgName {
			pkgPath = strings.Trim(imp.Path.Value, `"`)
			isAliasPkgName = true
			break
		}

		// 普通导入，包没有别名
		path := strings.Trim(imp.Path.Value, `"`)
		paths := strings.Split(path, "/")
		if paths[len(paths)-1] == pkgName {
			// 找到包路径
			pkgPath = path
			break
		}
	}
	return
}

func getFullTypeName(pkgName, typeName string) string {
	if pkgName != "" {
		return pkgName + "." + typeName
	}

	return typeName
}
