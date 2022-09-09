package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/whaios/goshowdoc/log"
)

var (
	listDoc   = `{"Title":"获取书籍列表","Catalog":"测试文档/书籍","Description":"分页获取书籍列表","Remark":"","Order":"1","Request":{"Method":"get","Url":"{{BASEURL}}/api/v1/book/list","ApiStatus":"","Headers":[{"name":"Authorization","type":"string","require":"1","value":"bearer {{TOKEN}}","remark":"用户登录凭证"}],"PathVariable":[],"Query":[{"name":"page","type":"int","require":"1","value":"","remark":"第几页"},{"name":"page_size","type":"int","require":"1","value":"","remark":"每页显示条数"}],"ParamMode":"urlencoded","Params":[],"ParamJson":""},"Response":{"Example":"{\n    \"errcode\": 0,\n    \"errmsg\": \"错误说明\",\n    \"data\": {\n        \"total_count\": 0,\n        \"items\": [\n            {\n                \"id\": \"标识符\",\n                \"title\": \"书名\",\n                \"publisher\": \"出版社\",\n                \"tags\": [\n                    \"\"\n                ]\n            }\n        ]\n    }\n}","Params":[{"name":"errcode","type":"int","remark":"错误代码"},{"name":"errmsg","type":"string","remark":"错误说明"},{"name":"total_count","type":"int","remark":"总条数"},{"name":"items","type":"array","remark":"书籍"},{"name":"items.id","type":"string","remark":"标识符"},{"name":"items.title","type":"string","remark":"书名"},{"name":"items.publisher","type":"string","remark":"出版社"},{"name":"items.tags","type":"array","remark":"标签"}]},"ResponseFail":{"Example":"","Params":[]}}`
	detailDoc = `{"Title":"获取指定书籍详情","Catalog":"测试文档/书籍","Description":"","Remark":"","Order":"2","Request":{"Method":"get","Url":"{{BASEURL}}/api/v1/book/detail/:id","ApiStatus":"","Headers":[{"name":"Authorization","type":"string","require":"1","value":"bearer {{TOKEN}}","remark":"用户登录凭证"}],"PathVariable":[{"name":"id","type":"int","require":"1","value":"","remark":"书籍 id"}],"Query":[],"ParamMode":"urlencoded","Params":[],"ParamJson":""},"Response":{"Example":"{\n    \"errcode\": 0,\n    \"errmsg\": \"错误说明\",\n    \"data\": {\n        \"id\": \"id\",\n        \"title\": \"书名\",\n        \"type\": \"包装：平装、精装\",\n        \"pages\": 0,\n        \"pub_date\": 0,\n        \"publisher\": \"出版社\",\n        \"isbn\": \"图书编号\",\n        \"is_active\": false,\n        \"desc\": \"介绍\",\n        \"pub_date_str\": \"出版日期\",\n        \"reviews\": [\n            {\n                \"id\": \"评论id\",\n                \"creation_unix\": 0,\n                \"book_id\": \"书籍id\",\n                \"content\": \"评论内容\",\n                \"review_user_id\": 0,\n                \"review_user_name\": \"评论人名称\",\n                \"recursive_reviews\": []\n            }\n        ],\n        \"review_page\": {\n            \"page\": 0,\n            \"page_size\": 0\n        }\n    }\n}","Params":[{"name":"errcode","type":"int","remark":"错误代码"},{"name":"errmsg","type":"string","remark":"错误说明"},{"name":"id","type":"string","remark":"id"},{"name":"title","type":"string","remark":"书名"},{"name":"type","type":"string","remark":"包装：平装、精装"},{"name":"pages","type":"int","remark":"页数"},{"name":"pub_date","type":"long","remark":"出版日期"},{"name":"publisher","type":"string","remark":"出版社"},{"name":"isbn","type":"string","remark":"图书编号"},{"name":"is_active","type":"boolean","remark":"是否激活"},{"name":"desc","type":"string","remark":"介绍"},{"name":"pub_date_str","type":"string","remark":"出版日期"},{"name":"reviews","type":"array","remark":"书籍评论"},{"name":"reviews.id","type":"object","remark":"评论id"},{"name":"reviews.creation_unix","type":"long","remark":"发表时间"},{"name":"reviews.book_id","type":"object","remark":"书籍id"},{"name":"reviews.content","type":"string","remark":"评论内容"},{"name":"reviews.review_user_id","type":"long","remark":"评论人id"},{"name":"reviews.review_user_name","type":"string","remark":"评论人名称"},{"name":"reviews.recursive_reviews","type":"array","remark":"测试是否能安全解析递归类型"},{"name":"review_page","type":"object","remark":"书籍评论分页"},{"name":"review_page.page","type":"int","remark":"第几页"},{"name":"review_page.page_size","type":"int","remark":"每页显示条数"}]},"ResponseFail":{"Example":"","Params":[]}}`
	editDoc   = `{"Title":"新建或编辑书籍","Catalog":"测试文档/书籍/管理","Description":"","Remark":"","Order":"3","Request":{"Method":"post","Url":"{{BASEURL}}/api/v1/book/edit","ApiStatus":"","Headers":[{"name":"Authorization","type":"string","require":"1","value":"bearer {{TOKEN}}","remark":"用户登录凭证"}],"PathVariable":[],"Query":[],"ParamMode":"json","Params":[{"name":"id","type":"string","require":"0","value":"","remark":"id"},{"name":"title","type":"string","require":"1","value":"","remark":"书名"},{"name":"type","type":"string","require":"0","value":"","remark":"包装：平装、精装"},{"name":"pages","type":"int","require":"0","value":"0","remark":"页数"},{"name":"pub_date","type":"long","require":"0","value":"0","remark":"出版日期"},{"name":"publisher","type":"string","require":"0","value":"","remark":"出版社"},{"name":"isbn","type":"string","require":"0","value":"","remark":"图书编号"},{"name":"is_active","type":"boolean","require":"0","value":"false","remark":"是否激活"}],"ParamJson":"{\n    \"id\": \"id\",\n    \"title\": \"书名\",\n    \"type\": \"包装：平装、精装\",\n    \"pages\": 0,\n    \"pub_date\": 0,\n    \"publisher\": \"出版社\",\n    \"isbn\": \"图书编号\",\n    \"is_active\": false\n}"},"Response":{"Example":"{\n    \"errcode\": 0,\n    \"errmsg\": \"错误说明\"\n}","Params":[{"name":"errcode","type":"int","remark":"错误代码"},{"name":"errmsg","type":"string","remark":"错误说明"}]},"ResponseFail":{"Example":"","Params":[]}}`
	delDoc    = `{"Title":"删除书籍","Catalog":"测试文档/书籍/管理","Description":"","Remark":"危险操作","Order":"4","Request":{"Method":"delete","Url":"{{BASEURL}}/api/v1/book/del/:id","ApiStatus":"","Headers":[{"name":"Authorization","type":"string","require":"1","value":"bearer {{TOKEN}}","remark":"用户登录凭证"}],"PathVariable":[{"name":"id","type":"int","require":"1","value":"","remark":"书籍 id"}],"Query":[],"ParamMode":"urlencoded","Params":[],"ParamJson":""},"Response":{"Example":"{\n    \"errcode\": 0,\n    \"errmsg\": \"错误说明\"\n}","Params":[{"name":"errcode","type":"int","remark":"错误代码"},{"name":"errmsg","type":"string","remark":"错误说明"}]},"ResponseFail":{"Example":"","Params":[]}}`
)

func TestParseApiDoc(t *testing.T) {
	Convey("测试解析注释", t, func() {
		dir := "../example/ginweb/handler"

		p := NewParser()
		So(p.ParseApiDoc(dir), ShouldBeNil)
		So(len(p.Docs), ShouldEqual, 4)

		wantDocs := []string{
			listDoc,
			detailDoc,
			editDoc,
			delDoc,
		}
		for i, doc := range p.Docs {
			So(doc.Json(), ShouldEqual, wantDocs[i])
		}
	})
}

func TestParseObject_ListRsp(t *testing.T) {
	Convey("测试解析对象", t, func() {
		log.IsDebug = true
		searchDir := "../example/ginweb/handler"
		typeName := "ginweb/handler/book.ListRsp"

		p := NewParser()
		So(p.collectGoFile(searchDir), ShouldBeNil)

		obj, err := p.ParseObject(typeName, nil) // 测试时 file == nil，所有需要typeName有包名
		So(err, ShouldBeNil)
		So(obj, ShouldNotBeNil)
		for _, f := range obj.AllFields() {
			Println("> "+typeName+":", f.Name, f.Type, f.Required, f.Value, f.Comment)
		}

		wantJson := `{"total_count":0,"items":[{"id":"标识符","title":"书名","publisher":"出版社","tags":[""]}]}`
		So(string(obj.Json()), ShouldEqual, wantJson)
	})
}

func TestParseObject_Detail(t *testing.T) {
	Convey("测试解析对象", t, func() {
		log.IsDebug = true
		searchDir := "../example/ginweb/handler"
		typeName := "ginweb/handler/book.Detail"

		p := NewParser()
		So(p.collectGoFile(searchDir), ShouldBeNil)

		obj, err := p.ParseObject(typeName, nil)
		So(err, ShouldBeNil)
		So(obj, ShouldNotBeNil)
		for _, f := range obj.AllFields() {
			Println("> "+typeName+":", f.Name, f.Type, f.Required, f.Value, f.Comment)
		}
		wantJson := `{"id":"id","title":"书名","type":"包装：平装、精装","pages":0,"pub_date":0,"publisher":"出版社","isbn":"图书编号","is_active":false,"desc":"介绍","pub_date_str":"出版日期","reviews":[{"id":"评论id","creation_unix":0,"book_id":"书籍id","content":"评论内容","review_user_id":0,"review_user_name":"评论人名称","recursive_reviews":[]}],"review_page":{"page":0,"page_size":0}}`
		So(string(obj.Json()), ShouldEqual, wantJson)
	})
}
