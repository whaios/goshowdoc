package parser

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/whaios/goshowdoc/log"
	"golang.org/x/tools/go/packages"
)

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
	projectDir string // 目标Go项目所在目录（支持绝对路径和相对路径），用于加载外部包时获取go包名，为空时默认为当前运行目录。

	files    map[*ast.File]*AstFileInfo // 使用到的所有go文件
	packages map[string]*Package        // key=完整包名. 如：ginweb/book
	// 在目录下收集到的具有唯一名称（完整包名+类型名）的类型。如果存在同名的不会出现到该字典中。
	uniqueDefinitions map[string]*TypeSpecDef // key=类型全名. 如：ginweb.handler.book.Book
}

// AddFile 添加 Go 源码文件，并解析代码中的类型
func (p *Packages) AddFile(pkgPath, fileName string, astFile *ast.File) *AstFileInfo {
	info := &AstFileInfo{
		File:     astFile,
		FileName: fileName,
		PkgPath:  pkgPath,
	}
	p.files[astFile] = info
	log.Debug("收集Go文件: %s", info.FileName)

	// 解析go文件中的结构体
	p.parseTypesFromFile(info.File, info.PkgPath)
	return info
}

// 解析go文件中的结构体
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
// @param shortName 包名.类型名，如：ListRsp 或 book.Book
func (p *Packages) FindTypeSpec(shortName string, file *ast.File) *TypeSpecDef {
	if file == nil { // for test
		return p.uniqueDefinitions[shortName]
	}

	var pkgName, typeName string
	{
		typeName = shortName
		if parts := strings.SplitN(shortName, ".", 2); len(parts) == 2 {
			pkgName = parts[0]
			typeName = parts[1]
		}
	}

	// 有包名，查找外部包
	if pkgName != "" {
		// 从文件中导入的包中查找指定包路径
		imptPkgPath, _ := p.findPackagePathFromImports(pkgName, file)
		// 没有找到对应的包名
		if imptPkgPath == "" {
			return nil
		}

		// 收集外部包
		p.loadExternalPackage(imptPkgPath)
		return p.findTypeSpec(imptPkgPath, typeName)
	}

	var pkgPath, fullName string
	if fileInfo, ok := p.files[file]; ok {
		pkgPath = fileInfo.PkgPath
		fullName = pkgPath + "." + typeName
	}

	// 从目录包中查找
	typeDef, ok := p.uniqueDefinitions[fullName]
	if ok {
		return typeDef
	}

	typeDef = p.findTypeSpec(pkgPath, typeName)
	if typeDef != nil {
		return typeDef
	}

	// 载入 . 包
	for _, imp := range file.Imports {
		if imp.Name != nil && imp.Name.Name == "." {
			imptPkgPath := strings.Trim(imp.Path.Value, `"`)
			// 收集外部包
			p.loadExternalPackage(imptPkgPath)

			if typeDef = p.findTypeSpec(imptPkgPath, typeName); typeDef != nil {
				return typeDef
			}
		}
	}
	return nil
}

// findTypeSpec 从收集的指定包中查找类型
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

// imptPkgPath 加载指定包
func (p *Packages) loadExternalPackage(imptPkgPath string) {
	if p.packages != nil {
		if _, ok := p.packages[imptPkgPath]; ok {
			// 已经收集过该包
			return
		}
	}
	log.Debug("加载外部包: %s", imptPkgPath)

	cfg := &packages.Config{
		Dir:  p.projectDir,
		Mode: packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedCompiledGoFiles,
	}
	pkgs, _ := packages.Load(cfg, imptPkgPath)
	for _, pkg := range pkgs {
		for i, astFile := range pkg.Syntax {
			p.AddFile(pkg.ID, pkg.CompiledGoFiles[i], astFile)
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

// AstFileInfo ast.File 文件信息.
type AstFileInfo struct {
	File     *ast.File
	FileName string // Go 源码文件全名称
	PkgPath  string // Go 源码文件完整包名
}
