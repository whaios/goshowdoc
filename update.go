package main

import (
	"github.com/whaios/goshowdoc/log"
	"github.com/whaios/goshowdoc/parser"
	"github.com/whaios/goshowdoc/runapi"
)

// Update 更新文档
func Update(searchDir string) {
	log.Info("解析Go源码文件 %s", searchDir)
	p := parser.NewParser()
	if err := p.ParseApiDoc(searchDir); err != nil {
		log.Error(err.Error())
		return
	}

	max := len(p.Docs)
	for i, doc := range p.Docs {
		if err := runapi.UpdateByApi(doc.Catalog, doc.Title, doc.Order, apiDocToPageContent(doc)); err != nil {
			log.Error("更新文档[%s/%s]失败: %s", doc.Catalog, doc.Title, err.Error())
			return
		}
		log.DrawProgressBar("更新文档", i+1, max)
	}
	log.Success("更新完成")
	return
}

func apiDocToPageContent(doc *parser.ApiDoc) *runapi.PageContent {
	content := runapi.NewPageContent(doc.Request.Method, doc.Request.Url)
	content.Info.Title = doc.Title
	content.Info.Description = doc.Description
	content.Info.Remark = doc.Remark

	if len(doc.Request.Headers) > 0 {
		content.Request.Headers = doc.Request.Headers
	}
	content.Request.Params.Mode = doc.Request.ParamMode
	if len(doc.Request.Params) > 0 {
		switch content.Request.Params.Mode {
		case runapi.ParamModeUrlEncoded:
			content.Request.Params.Urlencoded = doc.Request.Params
		case runapi.ParamModeFormData:
			content.Request.Params.Formdata = doc.Request.Params
		case runapi.ParamModeJson:
			content.Request.Params.JsonDesc = doc.Request.Params
		}
	}
	content.Request.Params.Json = doc.Request.ParamJson

	content.Response.ResponseExample = doc.Response.Example
	if len(doc.Response.Params) > 0 {
		content.Response.ResponseParamsDesc = doc.Response.Params
	}
	content.Response.Remark = doc.Remark
	return content
}
