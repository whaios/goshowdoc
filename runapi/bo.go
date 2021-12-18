package runapi

import (
	"encoding/json"
	"errors"
	"html"
	"strings"
)

type ErrResult struct {
	ErrorCode    int    `json:"error_code"` // 返回 0 表示成功
	ErrorMessage string `json:"error_message"`
}

func (p *ErrResult) Error() error {
	if p.ErrorCode != 0 {
		return errors.New(p.ErrorMessage)
	}
	return nil
}

// Item 项目
type Item struct {
	ItemId   string `json:"item_id"`
	ItemName string `json:"item_name"`
	Menu     struct {
		Catalogs []*Catalog  `json:"catalogs"`
		Pages    []*MenuPage `json:"pages"`
	} `json:"menu"`
}

// Catalogs 获取项目下所有的目录
func (p *Item) Catalogs() []*Catalog {
	return collCatalogs(p.Menu.Catalogs)
}

func collCatalogs(catalogs []*Catalog) []*Catalog {
	cats := catalogs
	for _, cat := range catalogs {
		if len(cat.Catalogs) > 0 {
			cats = append(cats, collCatalogs(cat.Catalogs)...)
		}
	}
	return cats
}

// MenuPages 获取项目下所有的接口文档
func (p *Item) MenuPages() []*MenuPage {
	pages := p.Menu.Pages
	pages = append(pages, collPages(p.Menu.Catalogs)...)
	return pages
}

func collPages(catalogs []*Catalog) []*MenuPage {
	pages := make([]*MenuPage, 0)
	for _, cat := range catalogs {
		if len(cat.Pages) > 0 {
			pages = append(pages, cat.Pages...)
		}
		if len(cat.Catalogs) > 0 {
			pages = append(pages, collPages(cat.Catalogs)...)
		}
	}
	return pages
}

// Catalog 目录
type Catalog struct {
	ItemId  string `json:"item_id"`
	CatId   string `json:"cat_id"`
	CatName string `json:"cat_name"`
	Level   string `json:"level"`

	ParentCatId string      `json:"parent_cat_id"`
	Catalogs    []*Catalog  `json:"catalogs"`
	Pages       []*MenuPage `json:"pages"`
}

// MenuPage 接口文档
type MenuPage struct {
	CatId     string `json:"cat_id"`
	PageId    string `json:"page_id"`
	PageTitle string `json:"page_title"`
}

// Page 接口文档
type Page struct {
	ItemId      string `json:"item_id"`
	CatId       string `json:"cat_id"`
	PageId      string `json:"page_id"`
	PageTitle   string `json:"page_title"`
	PageContent string `json:"page_content"`
}

// HtmlUnescape HTML解码
func (p *Page) HtmlUnescape() {
	p.PageContent = html.UnescapeString(p.PageContent)
}

// ToMap 转换为请求参数
func (p *Page) ToMap() map[string]string {
	m := make(map[string]string)

	p.HtmlUnescape()
	data, _ := json.Marshal(p)
	json.Unmarshal(data, &m)
	return m
}

const (
	MethodPost    = "POST"
	MethodGet     = "GET"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodHead    = "HEAD"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

const (
	ParamModeUrlEncoded = "urlencoded"
	ParamModeFormData   = "formdata"
	ParamModeJson       = "json"
)

func NewPageContent(method, url string) *PageContent {
	c := &PageContent{}
	c.Info.From = "runapi"
	c.Info.Type = "api"
	c.Request.Params.Mode = ParamModeUrlEncoded
	c.Request.Params.Urlencoded = []RequestParam{{Type: ParamTypeString, Require: "1"}}
	c.Request.Params.Formdata = []RequestParam{{Type: ParamTypeString, Require: "1"}}
	c.Request.Params.JsonDesc = []RequestParam{{Type: ParamTypeString, Require: "1"}}
	c.Request.Headers = []RequestParam{{Type: ParamTypeString, Require: "1"}}
	c.Request.Cookies = []NameValue{{}}
	c.Response.ResponseParamsDesc = []ResponseParam{{Type: ParamTypeString}}

	c.Info.Method = method
	c.Info.Url = url
	return c
}

// PageContent API 接口详情
type PageContent struct {
	Info struct {
		From        string `json:"from"`        // 文档来源：固定值 "runapi"
		Type        string `json:"type"`        // 类型：固定值 "api"
		Title       string `json:"title"`       // 接口标题
		Description string `json:"description"` // 可选。接口简要描述
		Method      string `json:"method"`      // 请求方式：POST GET PUT DELETE HEAD CONNECT OPTIONS TRACE
		Url         string `json:"url"`         // 请求 URL 地址
		Remark      string `json:"remark"`      // 可选，备注信息，会自动生成到文档的末尾。
	} `json:"info"` // 接口文档信息
	Request struct {
		Params struct {
			Mode       string         `json:"mode"` // 参数类型：urlencoded formdata json
			Urlencoded []RequestParam `json:"urlencoded"`
			Formdata   []RequestParam `json:"formdata"`
			Json       string         `json:"json"`
			JsonDesc   []RequestParam `json:"jsonDesc"`
		} `json:"params"` // 请求参数
		Headers []RequestParam `json:"headers"` // Headers
		Cookies []NameValue    `json:"cookies"` // Cookies
		Auth    []struct{}     `json:"auth"`
	} `json:"request"` // 请求内容
	Response struct {
		ResponseText     string   `json:"responseText"`
		ResponseOriginal string   `json:"responseOriginal"`
		ResponseHeader   struct{} `json:"responseHeader"`
		ResponseStatus   int      `json:"responseStatus"`

		ResponseExample    string          `json:"responseExample"`    // 返回示例
		ResponseParamsDesc []ResponseParam `json:"responseParamsDesc"` // 返回参数说明
		Remark             string          `json:"remark"`             // 可选，备注信息，会自动生成到文档的末尾。
	} `json:"response"` // 返回示例和参数说明
	Scripts struct {
		Pre  string `json:"pre"`  // 前执行脚本
		Post string `json:"post"` // 后执行脚本
	} `json:"scripts"` // 执行脚本
	Extend struct{} `json:"extend"`
}

func (p *PageContent) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

const (
	ParamTypeString  = "string"
	ParamTypeNumber  = "number"
	ParamTypeArray   = "array"
	ParamTypeObject  = "object"
	ParamTypeInt     = "int"
	ParamTypeLong    = "long"
	ParamTypeDate    = "date"
	ParamTypeBoolean = "boolean"
)

func transToHeaderType(typeName string) string {
	switch typeName {
	case "uint", "int", "uint8", "int8", "uint16", "int16", "byte",
		"uint32", "int32", "rune",
		"uint64", "int64",
		"float32", "float64":
		return ParamTypeNumber
	default:
		return ParamTypeString
	}
}

func transToDataType(typeName string) string {
	if strings.HasPrefix(typeName, "[]") {
		return ParamTypeArray
	}
	switch typeName {
	case "uint", "int", "uint8", "int8", "uint16", "int16", "byte":
		return ParamTypeInt
	case "uint32", "int32", "rune":
		return ParamTypeInt
	case "uint64", "int64":
		return ParamTypeLong
	case "float32", "float64":
		return ParamTypeNumber
	case "bool":
		return ParamTypeBoolean
	case "string":
		return ParamTypeString
	}

	return ParamTypeObject
}

func trasToRequire(require string) string {
	if require == "true" || require == "1" || require == "是" || require == "必填" {
		return "1"
	}
	return "0"
}

func NewHeaderParam(name, tpe, require, value, remark string) RequestParam {
	return RequestParam{
		Name:    name,
		Type:    transToHeaderType(tpe),
		Value:   value,
		Require: trasToRequire(require),
		Remark:  remark,
	}
}

func NewRequestParam(name, tpe, require, value, remark string) RequestParam {
	return RequestParam{
		Name:    name,
		Type:    transToDataType(tpe),
		Value:   value,
		Require: trasToRequire(require),
		Remark:  remark,
	}
}

// RequestParam API 接口请求参数
type RequestParam struct {
	Name    string `json:"name"`    // 字段名
	Type    string `json:"type"`    // 类型
	Require string `json:"require"` // 1=必填，0=选填
	Value   string `json:"value"`   // 值
	Remark  string `json:"remark"`  // 选填，header描述
}

type NameValue struct {
	Name  string `json:"name"`  // Cookie名
	Value string `json:"value"` // Cookie值
}

func NewResponseParam(name, tpe, remark string) ResponseParam {
	return ResponseParam{
		Name:   name,
		Type:   transToDataType(tpe),
		Remark: remark,
	}
}

// ResponseParam API 接口返回参数
type ResponseParam struct {
	Name   string `json:"name"`   // 字段名
	Type   string `json:"type"`   // 类型
	Remark string `json:"remark"` // 选填，字段描述。
}
