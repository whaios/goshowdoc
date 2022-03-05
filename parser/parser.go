package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/whaios/goshowdoc/log"
)

func NewParser() *Parser {
	return &Parser{
		Docs:     make([]*ApiDoc, 0),
		packages: NewPackages(),
	}
}

type Parser struct {
	packages *Packages // 注释中可能会引用到的类型
	Docs     []*ApiDoc // 解析注释生成的文档
}

// ParseApiDoc 解析指定目录下的 Go 代码文件注释，并生成文档。
// @param searchDir 目录下必须有 Go 代码文件
func (p *Parser) ParseApiDoc(searchDir string) error {
	if err := p.parseDir(searchDir); err != nil {
		return err
	}

	// 解析所有Go源码文件中的注释
	for _, fileInfo := range p.packages.SortedFiles() {
		if err := p.parseApiDoc(fileInfo.Path, fileInfo.File); err != nil {
			return err
		}
	}
	return nil
}

// parseApiDoc 将Go源码文件中的注释解析为API文档
func (p *Parser) parseApiDoc(fileName string, astFile *ast.File) error {
	var generalDoc = newApiDoc(p, astFile, nil)
	var order int64 = 1
	for _, astDescription := range astFile.Decls {
		switch astDescription.(type) {
		case *ast.GenDecl:
			// 获取类型上的通用注释
			astDecl := astDescription.(*ast.GenDecl)
			if astDecl.Doc != nil && astDecl.Doc.List != nil {
				log.Debug("解析通用注释: %s", fileName)
				for _, comment := range astDecl.Doc.List {
					if err := generalDoc.ParseComment("", comment.Text); err != nil {
						return fmt.Errorf("解析通用注释出错 %s :%+v", fileName, err)
					}
				}
			}
		case *ast.FuncDecl:
			// 获取方法上的注释
			astDecl := astDescription.(*ast.FuncDecl)
			if astDecl.Doc != nil && astDecl.Doc.List != nil {
				log.Debug("解析方法注释: %s %s()", fileName, astDecl.Name.Name)
				doc := newApiDoc(p, astFile, generalDoc)
				// 逐行解析方法上的注释块
				for _, comment := range astDecl.Doc.List {
					log.Debug("	> 注释: %s", comment.Text)
					if err := doc.ParseComment(astDecl.Name.Name, comment.Text); err != nil {
						return fmt.Errorf("解析方法注释出错 %s %s():%+v", fileName, astDecl.Name.Name, err)
					}
				}
				// 检查是否合法的API文档
				if doc.Invalid() {
					log.Debug("忽略方法注释（没有 title 或 url）: %s()", astDecl.Name.Name)
					continue
				}

				doc.Order = strconv.FormatInt(order, 10)
				log.Info("生成文档(%d) %s", order, doc.Name())
				p.Docs = append(p.Docs, doc)
				order++
			}
		}
	}
	return nil
}

// parseDir 收集指定目录下的 Go 代码文件和类型。
func (p *Parser) parseDir(searchDir string) error {
	p.packages.projectDir = searchDir

	// 收集指定目录下的 Go 代码文件
	if err := p.getAllGoFileInfo(searchDir); err != nil {
		return err
	}
	// 解析代码中的类型
	p.packages.ParseTypes()
	return nil
}

// 获取指定目录的包名："./example/ginweb/handler" > "ginweb/handler"
func getPkgName(searchDir string) (string, error) {
	cmd := exec.Command("go", "list", "-f={{.ImportPath}}")
	cmd.Dir = searchDir
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execute go list command, %s, stdout:%s, stderr:%s", err, stdout.String(), stderr.String())
	}

	outStr, _ := stdout.String(), stderr.String()

	if outStr[0] == '_' { // will shown like _/{GOPATH}/src/{YOUR_PACKAGE} when NOT enable GO MODULE.
		outStr = strings.TrimPrefix(outStr, "_"+build.Default.GOPATH+"/src/")
	}
	f := strings.Split(outStr, "\n")
	outStr = f[0]

	return outStr, nil
}

// GetAllGoFileInfo 获取指定目录下的所有Go代码文件.
// @param searchDir 如：./example/ginweb/handler
func (p *Parser) getAllGoFileInfo(searchDir string) error {
	// packageDir = ginweb/handler
	packageDir, err := getPkgName(searchDir)
	if err != nil {
		return fmt.Errorf("获取包名失败, dir: %s, error: %s", searchDir, err.Error())
	}

	return filepath.Walk(searchDir, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), "_test.go") || filepath.Ext(path) != ".go" {
			return nil
		}

		// path = example\ginweb\handler\handler.go
		// filePkg = ginweb/handler
		var filePkg string
		{
			relPath, err := filepath.Rel(searchDir, path)
			if err != nil {
				return err
			}
			filePkg = filepath.ToSlash(filepath.Dir(filepath.Clean(filepath.Join(packageDir, relPath))))
		}

		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("ParseFile error:%+v", err)
		}

		// absPath = D:\Work\github.com\whaios\goshowdoc\example\ginweb\handler\handler.go
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		p.packages.CollectAstFile(filePkg, absPath, astFile)
		return nil
	})
}

// ParseObject 解析对象
func (p *Parser) ParseObject(typeName string, file *ast.File) (*Object, error) {
	log.Debug("解析类型: %s", typeName)
	typeSpecDef := p.packages.FindTypeSpec(typeName, file)
	if typeSpecDef == nil {
		return nil, fmt.Errorf("没有找到类型定义: %s", typeName)
	}

	// 解析结构体字段
	st, ok := typeSpecDef.TypeSpec.Type.(*ast.StructType)
	if !ok {
		// 不是有效的类型，可能是自定义基础类型
		return nil, nil
	}

	obj := New()
	for _, field := range st.Fields.List {
		var name, jsonName, dataType, tag, comment string
		var required bool

		dataType = parseFieldType(field.Type)

		if len(field.Names) == 0 {
			// 匿名字段
			nObj, err := p.ParseObject(dataType, typeSpecDef.File)
			if err != nil {
				return nil, err
			}
			obj.PutAnonymousObject(nObj)
			continue
		} else {
			name = field.Names[0].Name
		}

		if field.Tag != nil {
			// `json:"name" validate:"required"`
			tag = field.Tag.Value
			required = strings.Contains(tag, "required")

			if jsonTag := getJsonTag(tag); jsonTag != "" {
				var tagOpts tagOptions
				jsonName, tagOpts = parseJsonTag(jsonTag)
				if tagOpts.Contains("string") {
					// json 标签中定义了类型转换
					dataType = "string"
				}
			}

			if jsonName == "" {
				jsonName = name
			}
		}
		if field.Doc != nil {
			// 字段上行的注释
			// 忽略这种注释
		}
		if field.Comment != nil {
			// 字段后面的同行注释
			for _, comm := range field.Comment.List {
				comment += strings.TrimSpace(strings.TrimLeft(comm.Text, "//"))
			}
		}

		objField := NewField(jsonName, dataType, required, comment)

		if isGolangPrimitiveType(dataType) ||
			strings.HasPrefix(dataType, "map[") ||
			dataType == "interface{}" ||
			dataType == "" {
			// 基础数据类型字段
			obj.PutField(objField)
		} else if strings.HasPrefix(dataType, "[]") {
			// 切片类型字段
			itemType := strings.TrimLeft(dataType, "[]")
			if isGolangPrimitiveType(itemType) {
				obj.PutArray(objField)
			} else if itemType == typeSpecDef.TypeSpec.Name.Name {
				// 避免无限循环解析递归类型
				obj.PutArray(objField)
			} else {
				nObj, err := p.ParseObject(itemType, typeSpecDef.File)
				if err != nil {
					return nil, err
				}
				obj.PutObjectArray(objField, nObj)
			}
		} else {
			nObj, err := p.ParseObject(dataType, typeSpecDef.File)
			if err != nil {
				return nil, err
			}
			if nObj == nil {
				// 没有解析为有效类型，按普通字段处理
				obj.PutField(objField)
			} else {
				obj.PutObject(objField, nObj)
			}
		}
	}
	return obj, nil
}

func isGolangPrimitiveType(typeName string) bool {
	switch typeName {
	case "uint",
		"int",
		"uint8",
		"int8",
		"uint16",
		"int16",
		"byte",
		"uint32",
		"int32",
		"rune",
		"uint64",
		"int64",
		"float32",
		"float64",
		"bool",
		"string":
		return true
	}

	return false
}

func parseFieldType(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch expr.(type) {
	case *ast.Ident:
		id := expr.(*ast.Ident)
		if id.Obj != nil && id.Obj.Decl != nil {
			if ts, ok := id.Obj.Decl.(*ast.TypeSpec); ok {
				// 自定义类型（可能是基础类型，也可能是struct，struct返回空字符串）
				if _, ok = ts.Type.(*ast.Ident); ok {
					return parseFieldType(ts.Type)
				}
			}
		}
		return id.Name
	case *ast.ArrayType:
		arrt := expr.(*ast.ArrayType)
		return "[]" + parseFieldType(arrt.Elt)
	case *ast.MapType:
		mpt := expr.(*ast.MapType)
		kn := parseFieldType(mpt.Key)
		vn := parseFieldType(mpt.Value)
		return fmt.Sprintf("map[%s]%s", kn, vn)
	case *ast.SelectorExpr: // 包名.类型
		selt := expr.(*ast.SelectorExpr)
		pkgName := parseFieldType(selt.X)
		if pkgName == "" {
			return selt.Sel.Name
		}
		return fmt.Sprintf("%s.%s", pkgName, selt.Sel.Name)
	case *ast.StarExpr: // 指针
		star := expr.(*ast.StarExpr)
		return parseFieldType(star.X)
	case *ast.InterfaceType:
		return "interface{}"
	}
	return ""
}
