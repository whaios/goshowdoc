package parser

import (
	"testing"

	gen "github.com/darjun/json-gen"
)

func TestJsonGen(t *testing.T) {
	wantJson := `{"fstring":"string value","fint":2,"fbool":true,"ffloat":160,"farray":["arr1"],"child":{"cfstring":"child string value"}}`

	m := gen.NewMap()
	m.PutString("fstring", "string value")
	m.PutInt("fint", 2)
	m.PutBool("fbool", true)
	m.PutFloat("ffloat", 160)

	arr := gen.NewArray()
	arr.AppendString("arr1")
	m.PutArray("farray", arr)

	child := gen.NewMap()
	child.PutString("cfstring", "child string value")
	m.PutMap("child", child)

	gotJson := string(m.Serialize(nil))
	if gotJson != wantJson {
		t.Errorf("获取的JSON字符串不匹配：\n%s", gotJson)
	}
}
