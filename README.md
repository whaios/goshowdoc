# goshowdoc

[ShowDoc](https://github.com/star7th/showdoc) API 接口文档工具

[ShowDoc](www.showdoc.com.cn) 是一个非常适合IT团队的在线API文档、技术文档工具。搭配的RunApi客户端中一边调试接口、一边自动生成文档。

本工具可通过解析 Go 代码文件中的注释自动生成API文档，并且该文档支持通过RunApi客户端进行调试。

## 安装

1. 使用如下命令下载goshowdoc：

```shell
$ go get -u github.com/whaios/goshowdoc

# 1.16 及以上版本
$ go install github.com/whaios/goshowdoc@latest
```

## 命令说明

```
$ goshowdoc.exe

NAME:
   goshowdoc - ShowDoc API 接口文档工具

USAGE:
   goshowdoc.exe [global options] command [command options] [arguments...]

VERSION:
   0.1.0

DESCRIPTION:
   项目地址： https://github.com/whaios/goshowdoc
   支持以下功能：
   1. 通过代码注释生成 API 接口文档。
   2. 导出和导入 ShowDoc 项目。

COMMANDS:
   flags      查询应用全局相关参数。
   item       ShowDoc 项目导出和导入。
   update, u  解析 Go 源码中的注释，生成并更新 ShowDoc 文档。
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value   ShowDoc 地址。 (default: "https://www.showdoc.com.cn") [%GOSHOWDOC_HOST%]
   --debug        开启调试模式。 (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## 注释生成文档

### 注释格式

#### 通用API注释

| 注释    | 说明    | 示例    |
| ----------------------- | ----------------------- | ----------------------- |
| @catalog | 文档目录，多级目录用 `/` 隔开 | // @catalog 一级/二级/三级 |
| @header | 可选，请求头。格式为 `[字段名] [类型] [必填] ["值"] ["备注"]` | // @header Authorization string true "abc" "用户登录凭证" |
| @response, @resp | 返回内容，支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] ["备注"]`）两种方式。 | // @response TestApiRsp{}  // @param page int "第几页" |
| @remark | 可选，备注信息 | // @remark 用户需要先登录 |

#### API 注释

| 注释           | 说明         | 示例            |
| ------------- | ----------- | -------------- |
| @title                | 接口文档标题，方法注释。 | // funcName 获取书籍列表 // @title 获取书籍列表  |
| @catalog              | 文档目录，多级目录用 `/` 隔开 | // @catalog 一级/二级/三级 |
| @url                  | 接口URL，格式为：`[method] [url]` | // @url GET {{BASEURL}}/api/v1/book/list |
| @description, @desc   | 可选，接口描述信息 | // @description 分页获取书籍列表 |
| @header               | 可选，请求头。格式为 `[字段名] [类型] [必填] ["值"] ["备注"]` | // @header Authorization string true "abc" "用户登录凭证" |
| @param_mode           | 可选，请求参数方式。`urlencoded(GET请求默认)`、`json(POST请求默认)` 和 `formdata` | // @param_mode urlencoded |
| @param                | 可选，请求参数。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] [必填] ["值"] ["备注"]`）两种方式。 | // @param id int true "" "书籍 id" |
| @response, @resp      | 可选，返回内容。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] ["备注"]`）两种方式。 | // @response TestApiRsp{}  // @param page int "第几页" |
| @remark               | 可选，备注信息 | // @remark 用户需要先登录 |

#### 代码示例

解析 Go 源码中的注释，生成并更新 ShowDoc 文档。

```shell
# 切换到项目目录
# cd example/ginweb

# 生成并更新文档
goshowdoc.exe u --dir ./handler/

# 生成并更新文档，输出调试信息
goshowdoc.exe --debug u --dir ./handler/
```

```go
package handler


// Handler 书籍管理
//
// @catalog 测试文档/书籍
// @header Authorization string true "bearer {{TOKEN}}" "用户登录凭证"
// @resp comm.HttpCode{}
type Handler struct {
}

// List 获取书籍列表
//
// @description 分页获取书籍列表
// @url GET {{BASEURL}}/api/v1/book/list
// @param	page		int	true	""	"第几页"
// @param	page_size	int	true	""	"每页显示条数"
// @response ListRsp{}
func (h *Handler) List() {
}

// Detail 获取指定书籍详情
//
// @url GET {{BASEURL}}/api/v1/book/detail/:id
// @param :id int true "" "书籍 id"
// @response Detail{}
func (h *Handler) Detail() {
}

// CreateOrUpdate 新建或编辑书籍
//
// @catalog 管理
// @url POST {{BASEURL}}/api/v1/book/edit
// @param book.Book{}
func (h *Handler) CreateOrUpdate() {
}

// Delete
//
// @catalog 管理
// @title 删除书籍
// @url DELETE {{BASEURL}}/api/v1/book/del/:id
// @param :id int true "" "书籍 id"
// @remark 危险操作
func (h *Handler) Delete() {
}
```

## ShowDoc 项目导出导入

项目导入导出功能使用 ShowDoc 内部接口而非开放接口，所以可能会出现ShowDoc内部接口调整导致导入功能无法正常使用。

目前兼容 ShowDoc 版本为 **v2.9.14**