package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetJsonName(t *testing.T) {
	cases := []struct {
		Tag      string
		JsonName string
	}{
		{"", ""},
		{`validate:"required"`, ""},
		{`json:""`, ""},
		{`json:",string"`, ",string"},
		{`json:"name"`, "name"},
		{`json:"name,string"`, "name,string"},
		{`json:"name" validate:"required"`, "name"},
		{`json:"name,string" validate:"required"`, "name,string"},
	}

	Convey("测试获取json标签内容", t, func() {
		for _, c := range cases {
			name := getJsonTag(c.Tag)
			So(name, ShouldEqual, c.JsonName)
		}
	})
}

func TestTagParsing(t *testing.T) {
	name, opts := parseJsonTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}
