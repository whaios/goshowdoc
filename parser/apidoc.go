package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/tidwall/sjson"
	"github.com/whaios/goshowdoc/runapi"
)

func newApiDoc(parser *Parser, astFile *ast.File, generalDoc *ApiDoc) *ApiDoc {
	doc := &ApiDoc{
		parser:  parser,
		astFile: astFile,
		Order:   "99",
	}
	doc.Request.Headers = make([]runapi.RequestParam, 0)
	doc.Request.PathVariable = make([]runapi.RequestParam, 0)
	doc.Request.Query = make([]runapi.RequestParam, 0)
	doc.Request.Params = make([]runapi.RequestParam, 0)
	doc.Response.Params = make([]runapi.ResponseParam, 0)
	doc.ResponseFail.Params = make([]runapi.ResponseParam, 0)
	if generalDoc != nil {
		doc.Catalog = generalDoc.Catalog
		doc.Remark = generalDoc.Remark
		for _, header := range generalDoc.Request.Headers {
			doc.Request.Headers = append(doc.Request.Headers, header)
		}
		for _, param := range generalDoc.Response.Params {
			doc.Response.Params = append(doc.Response.Params, param)
		}
		doc.Response.Example = generalDoc.Response.Example
	}
	return doc
}

// ApiDoc API 接口文档
type ApiDoc struct {
	parser  *Parser
	astFile *ast.File

	Title       string
	Catalog     string // 例如 “一层/二层/三层”
	Description string
	Remark      string
	Order       string // 文档排序，默认 99

	Request      ApiRequest
	Response     ApiResponse
	ResponseFail ApiResponse
}

type ApiRequest struct {
	Method       string
	Url          string
	ApiStatus    string // 接口状态
	Headers      []runapi.RequestParam
	PathVariable []runapi.RequestParam // 路径参数
	Query        []runapi.RequestParam // GET 请求建议仅用 Query 参数
	ParamMode    string                // 参数类型：urlencoded formdata json
	Params       []runapi.RequestParam
	ParamJson    string
}

type ApiResponse struct {
	Example string
	Params  []runapi.ResponseParam
}

// ParseComment 解析单行注释
func (p *ApiDoc) ParseComment(funcName, comment string) error {
	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if commentLine == "" {
		// 没有注释内容
		return nil
	}

	funcName = strings.ToLower(funcName)
	attribute := strings.ToLower(strings.Fields(commentLine)[0])
	lineRemainder := strings.TrimSpace(commentLine[len(attribute):])

	var err error
	switch attribute {
	case funcName, "@title":
		p.Title = lineRemainder
	case "@catalog":
		p.parseCatalogComment(lineRemainder)
	case "@desc", "@description":
		p.parseDescriptionComment(lineRemainder)
	case "@url":
		err = p.parseUrlComment(lineRemainder)
	case "@api_status":
		err = p.parseApiStatusComment(lineRemainder)
	case "@header":
		err = p.parseHeaderComment(lineRemainder)
	case "@path_var":
		err = p.parsePathVarComment(lineRemainder)
	case "@query":
		err = p.parseQueryComment(lineRemainder)
	case "@param_mode":
		err = p.parseParamModeComment(lineRemainder)
	case "@param":
		err = p.parseParamComment(lineRemainder)
	case "@resp", "@response":
		err = p.parseResponseComment(lineRemainder)
	case "@resp_fail", "@response_fail":
		err = p.parseResponseFailComment(lineRemainder)
	case "@remark":
		p.parseRemarkComment(lineRemainder)
	}
	return err
}

func (p *ApiDoc) parseCatalogComment(commentLine string) {
	if !strings.HasPrefix(commentLine, "/") {
		commentLine = "/" + commentLine
	}
	p.Catalog += commentLine
	p.Catalog = strings.Trim(p.Catalog, "/")
}

// parseDescriptionComment 解析多行描述
func (p *ApiDoc) parseDescriptionComment(commentLine string) {
	if p.Description != "" {
		p.Description += "\n"
	}
	p.Description += commentLine
}

func (p *ApiDoc) parseRemarkComment(commentLine string) {
	if p.Remark != "" {
		p.Remark += "\n"
	}
	p.Remark += commentLine
}

// parseUrlComment 解析url。
// 如：get {{BASEURL}}/api/v1/book/list
func (p *ApiDoc) parseUrlComment(commentLine string) error {
	fields := strings.Fields(commentLine)
	if len(fields) != 2 {
		return fmt.Errorf("无法解析 url 注释 \"%s\"", commentLine)
	}
	p.Request.Method = strings.ToLower(fields[0]) // 和runapi保持一致使用小写
	p.Request.Url = fields[1]
	switch p.Request.Method {
	case runapi.MethodPost:
		p.Request.ParamMode = runapi.ParamModeJson
	default:
		p.Request.ParamMode = runapi.ParamModeUrlEncoded
	}
	return nil
}

// parseApiStatusComment 解析接口状态
func (p *ApiDoc) parseApiStatusComment(commentLine string) error {
	p.Request.ApiStatus = runapi.TransToApiStatus(commentLine)
	return nil
}

var reqParamPattern = regexp.MustCompile(`(\S+)[\s]+([\w]+)[\s]+([\w]+)[\s]+"([^"]*)"[\s]+"([^"]*)"`)

// parseHeaderComment 解析Header。
//
// 如：	page		int		true	"1"		"第几页"
//		[字段名]		[类型]	[必填]	[值]	[备注]
func (p *ApiDoc) parseHeaderComment(commentLine string) error {
	matches := reqParamPattern.FindStringSubmatch(commentLine)
	if len(matches) != 6 {
		return fmt.Errorf("无法解析 header 注释 \"%s\"\n不符合格式 [字段名] [类型] [必填] [\"值\"] [\"备注\"]", commentLine)
	}

	param := runapi.NewHeaderParam(matches[1], matches[2], matches[3], matches[4], matches[5])
	p.Request.Headers = append(p.Request.Headers, param)
	return nil
}

// parsePathVarComment 解析路径参数
func (p *ApiDoc) parsePathVarComment(commentLine string) error {
	matches := reqParamPattern.FindStringSubmatch(commentLine)
	if len(matches) != 6 {
		return fmt.Errorf("无法解析 path_var 注释 \"%s\"\n不符合格式 [字段名] [类型] [必填] [\"值\"] [\"备注\"]", commentLine)
	}

	param := runapi.NewRequestParam(matches[1], matches[2], matches[3], matches[4], matches[5])
	p.Request.PathVariable = append(p.Request.PathVariable, param)
	return nil
}

// parseQueryComment 解析Query参数（GET请求建议仅用Query参数）
func (p *ApiDoc) parseQueryComment(commentLine string) error {
	params, _, err := p.parseRequestParam(commentLine)
	if err != nil {
		return err
	}
	p.Request.Query = append(p.Request.Query, params...)
	return nil
}

// parseParamModeComment 解析请求参数模式
func (p *ApiDoc) parseParamModeComment(commentLine string) error {
	switch commentLine {
	case runapi.ParamModeUrlEncoded:
	case runapi.ParamModeFormData, runapi.ParamModeJson:
		if p.Request.Method == runapi.MethodGet {
			fmt.Errorf("GET 请求只支持 %s 模式", runapi.ParamModeUrlEncoded)
		}
	default:
		return fmt.Errorf("不支持 %s 请求参数模式", commentLine)
	}
	p.Request.ParamMode = commentLine
	return nil
}

// parseParamComment 解析Body参数
func (p *ApiDoc) parseParamComment(commentLine string) error {
	params, paramJson, err := p.parseRequestParam(commentLine)
	if err != nil {
		return err
	}
	if p.Request.Method == runapi.MethodGet {
		p.Request.Query = append(p.Request.Query, params...)
	} else {
		p.Request.Params = append(p.Request.Params, params...)
	}
	if p.Request.ParamMode == runapi.ParamModeJson {
		p.Request.ParamJson = paramJson
	}
	return nil
}

// parseParamComment 解析请求参数
//
// 如：	page		int		true	"1"		"第几页"
//		[字段名]		[类型]	[必填]	[值]	[备注]
func (p *ApiDoc) parseRequestParam(commentLine string) (params []runapi.RequestParam, paramJson string, err error) {
	if !strings.HasSuffix(commentLine, "{}") {
		matches := reqParamPattern.FindStringSubmatch(commentLine)
		if len(matches) != 6 {
			err = fmt.Errorf("无法解析 param 注释 \"%s\"\n不符合格式 [字段名] [类型] [必填] [\"值\"] [\"备注\"]", commentLine)
			return
		}

		param := runapi.NewRequestParam(matches[1], matches[2], matches[3], matches[4], matches[5])
		params = append(params, param)
		return
	}

	// 解析对象
	refType := strings.TrimRight(commentLine, "{}")
	if refType == "" || p.parser == nil {
		return
	}
	obj, err := p.parser.ParseObject(refType, p.astFile)
	if err != nil || obj == nil {
		return
	}

	requireVal := func(required bool) string {
		if required {
			return "1"
		}
		return "0"
	}
	for _, field := range obj.AllFields() {
		param := runapi.NewRequestParam(field.Name, field.Type, requireVal(field.Required), field.Value, field.Comment)
		params = append(params, param)
	}
	paramJson = jsonFormat(obj.Json())
	return
}

var respParamPattern = regexp.MustCompile(`(\S+)[\s]+([\w]+)[\s]+"([^"]*)"`)

// parseResponseComment 解析返回样例
//
// 如：	page		int		"第几页"
//		[字段名]		[类型]	[备注]
func (p *ApiDoc) parseResponseComment(commentLine string) error {
	params, paramJson, err := p.parseResponseParam(commentLine)
	if err != nil {
		return err
	}

	p.Response.Params = append(p.Response.Params, params...)
	if p.Response.Example == "" {
		p.Response.Example = jsonFormat(paramJson)
	} else if strings.HasPrefix(p.Response.Example, "{") {
		val, _ := sjson.SetRawBytes([]byte(p.Response.Example), "data", paramJson)
		p.Response.Example = jsonFormat(val)
	}
	return nil
}

// parseResponseFailComment 解析失败返回示例
func (p *ApiDoc) parseResponseFailComment(commentLine string) error {
	params, paramJson, err := p.parseResponseParam(commentLine)
	if err != nil {
		return err
	}

	p.ResponseFail.Params = append(p.ResponseFail.Params, params...)
	if p.ResponseFail.Example == "" {
		p.ResponseFail.Example = jsonFormat(paramJson)
	} else if strings.HasPrefix(p.ResponseFail.Example, "{") {
		val, _ := sjson.SetRawBytes([]byte(p.ResponseFail.Example), "data", paramJson)
		p.ResponseFail.Example = jsonFormat(val)
	}
	return nil
}

func (p *ApiDoc) parseResponseParam(commentLine string) (params []runapi.ResponseParam, paramJson []byte, err error) {
	if !strings.HasSuffix(commentLine, "{}") {
		matches := respParamPattern.FindStringSubmatch(commentLine)
		if len(matches) != 4 {
			err = fmt.Errorf("无法解析 response 注释 \"%s\"\n不符合格式 [字段名] [类型] [\"备注\"]", commentLine)
			return
		}
		param := runapi.ResponseParam{
			Name:   matches[1],
			Type:   matches[2],
			Remark: matches[3],
		}
		params = append(params, param)
		return
	}

	// 解析对象
	refType := strings.TrimRight(commentLine, "{}")
	if refType == "" || p.parser == nil {
		return
	}
	obj, err := p.parser.ParseObject(refType, p.astFile)
	if err != nil || obj == nil {
		return
	}

	for _, field := range obj.AllFields() {
		param := runapi.NewResponseParam(field.Name, field.Type, field.Comment)
		params = append(params, param)
	}
	paramJson = obj.Json()
	return
}

// Invalid 没有标题或Url，不是有效的API文档
func (p *ApiDoc) Invalid() bool {
	return p.Title == "" || p.Request.Url == ""
}

// Name 文档分类+标题
func (p *ApiDoc) Name() string {
	catalog := p.Catalog
	if catalog != "" {
		catalog += "/"
	}
	return catalog + p.Title
}

// Json 测试用
func (p *ApiDoc) Json() string {
	data, _ := json.Marshal(p)
	return string(data)
}

// 格式化JSON字符串
func jsonFormat(data []byte) string {
	var buf bytes.Buffer
	_ = json.Indent(&buf, data, "", "    ")
	return buf.String()
}
