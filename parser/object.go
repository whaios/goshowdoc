package parser

import (
	"fmt"
	"strings"

	gen "github.com/darjun/json-gen"
)

func New() *Object {
	return &Object{
		json:   gen.NewMap(),
		Fields: make([]*Field, 0),
	}
}

func NewField(name, typeName string, required bool, comment string) *Field {
	return &Field{
		Name:     name,
		Type:     typeName,
		Required: required,
		Comment:  comment,
	}
}

// Object 模拟对象
type Object struct {
	json *gen.Map

	Fields []*Field
}

// Field 字段属性
type Field struct {
	Name     string // 字段json名称
	Type     string // 字段类型
	Required bool   // 是否必填，字段tag中有required标记
	Value    string // 字段的模拟值
	Comment  string // 字段同行注释

	fields []*Field
}

// AllFields 所有字段数组，包含子对象字段
func (obj *Object) AllFields() []*Field {
	return getFields("", obj.Fields)
}

// Json 对象的json字符串
func (obj *Object) Json() []byte {
	return obj.json.Serialize(nil)
}

func getFields(parentName string, fields []*Field) []*Field {
	fs := make([]*Field, 0)
	for _, f := range fields {
		if parentName != "" {
			f.Name = fmt.Sprintf("%s.%s", parentName, f.Name)
		}
		fs = append(fs, f)
		if len(f.fields) > 0 {
			fs = append(fs, getFields(f.Name, f.fields)...)
		}
	}
	return fs
}

// PutField 添加基础类型字段
func (obj *Object) PutField(field *Field) {
	fv := ""
	switch field.Type {
	case "uint",
		"int",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int8",
		"int16",
		"int32",
		"int64":
		fv = "0"
		obj.json.PutInt(field.Name, 0)
	case "float32",
		"float64":
		fv = "0.00"
		obj.json.PutFloat(field.Name, 0.00)
	case "bool":
		fv = "false"
		obj.json.PutBool(field.Name, false)
	default:
		fv = ""
		obj.json.PutString(field.Name, field.Comment)
	}

	field.Value = fv
	obj.Fields = append(obj.Fields, field)
}

// PutAnonymousObject 添加匿名字段
func (obj *Object) PutAnonymousObject(value *Object) {
	for _, f := range value.Fields {
		obj.PutField(f)
	}
}

// PutArray 添加数组字段
func (obj *Object) PutArray(field *Field) {
	arr := gen.NewArray()

	tpe := strings.TrimLeft(field.Type, "[]")
	switch tpe {
	case "uint",
		"int",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int8",
		"int16",
		"int32",
		"int64":
		arr.AppendInt(0)
	case "float32",
		"float64":
		arr.AppendFloat(0.00)
	case "string":
		arr.AppendString("")
	}

	obj.json.PutArray(field.Name, arr)
	obj.Fields = append(obj.Fields, field)
}

// PutObjectArray 添加对象数组字段
func (obj *Object) PutObjectArray(field *Field, value *Object) {
	field.fields = value.Fields
	obj.Fields = append(obj.Fields, field)

	arr := gen.NewArray()
	arr.AppendMap(value.json)
	obj.json.PutArray(field.Name, arr)
}

// PutObject 添加对象字段
func (obj *Object) PutObject(field *Field, value *Object) {
	field.fields = value.Fields
	obj.Fields = append(obj.Fields, field)
	obj.json.PutMap(field.Name, value.json)
}
