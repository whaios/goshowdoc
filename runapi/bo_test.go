package runapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPageContent_SetQuery(t *testing.T) {
	params := []RequestParam{
		{Name: "q1", Value: "v1"},
		{Name: "q2", Value: ""},
	}
	cases := []struct {
		Url  string
		Want string
	}{
		{"http://url.com", "http://url.com?q1=v1&q2="},
		{"http://url.com?", "http://url.com?q1=v1&q2="},
		{"http://url.com?id=1", "http://url.com?id=1&q1=v1&q2="},
	}

	Convey("测试解析返回参数", t, func() {
		content := NewPageContent("GET", "")
		for _, cs := range cases {
			content.Info.Url = cs.Url
			content.SetQuery(params)
			So(content.Info.Url, ShouldEqual, cs.Want)
		}
	})
}
