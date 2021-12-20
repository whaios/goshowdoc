package parser

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/whaios/goshowdoc/runapi"
)

func TestApiDoc_ParseGeneralComment(t *testing.T) {
	Convey("测试解析通用注释", t, func() {
		funcName := "TestApi"
		title := "接口标题"
		catalog := "一级目录/二级目录"
		description := "可选 接口简要描述"
		remark := "可选 文档的末尾的备注信息"
		comments := []string{
			fmt.Sprintf("// %s %s", funcName, title),
			"// @catalog " + catalog,
			"// @description " + description,
			"// @remark " + remark,
		}

		doc := &ApiDoc{}
		for _, comment := range comments {
			So(doc.ParseComment(funcName, comment), ShouldBeNil)
		}
		So(doc.Title, ShouldEqual, title)
		So(doc.Catalog, ShouldEqual, catalog)
		So(doc.Description, ShouldEqual, description)
		So(doc.Remark, ShouldEqual, remark)
	})
}

func TestApiDoc_ParseUrlComment(t *testing.T) {
	Convey("测试解析 @url 注释", t, func() {
		url := "{{BASEURL}}/api/v1/book"

		doc := &ApiDoc{}
		So(doc.ParseUrlComment("GET "+url), ShouldBeNil)
		So(doc.Request.Url, ShouldEqual, url)
		So(doc.Request.Method, ShouldEqual, runapi.MethodGet)
		So(doc.Request.ParamMode, ShouldEqual, runapi.ParamModeUrlEncoded)

		So(doc.ParseUrlComment("POST "+url), ShouldBeNil)
		So(doc.Request.Url, ShouldEqual, url)
		So(doc.Request.Method, ShouldEqual, runapi.MethodPost)
		So(doc.Request.ParamMode, ShouldEqual, runapi.ParamModeJson)
	})
}

func TestApiDoc_ParseReqParamComment(t *testing.T) {
	Convey("测试解析请求参数", t, func() {
		paramComment := `page_size	int	true	"30"	"每页显示条数"`

		doc := &ApiDoc{}
		So(doc.ParseParamComment(paramComment), ShouldBeNil)
		So(len(doc.Request.Params) > 0, ShouldBeTrue)
		param := doc.Request.Params[0]
		So(param.Name, ShouldEqual, "page_size")
		So(param.Type, ShouldEqual, "int")
		So(param.Require, ShouldEqual, "1")
		So(param.Value, ShouldEqual, "30")
		So(param.Remark, ShouldEqual, "每页显示条数")

		paramComment = `page	int	true	""	""`
		So(doc.ParseParamComment(paramComment), ShouldBeNil)
		So(len(doc.Request.Params) > 1, ShouldBeTrue)
		So(doc.Request.Params[1].Name, ShouldEqual, "page")
	})
}

func TestApiDoc_ParseResponseParamComment(t *testing.T) {
	Convey("测试解析返回参数", t, func() {
		paramComment := `book.title string "书名"`

		doc := &ApiDoc{}
		So(doc.ParseResponseComment(paramComment), ShouldBeNil)
		So(len(doc.Response.Params), ShouldEqual, 1)
		param := doc.Response.Params[0]
		So(param.Name, ShouldEqual, "book.title")
		So(param.Type, ShouldEqual, "string")
		So(param.Remark, ShouldEqual, "书名")

		paramComment = `book.title string ""`
		So(doc.ParseResponseComment(paramComment), ShouldBeNil)
	})
}
