# goshowdoc

ShowDoc API 接口文档工具。

[ShowDoc](https://github.com/star7th/showdoc) 是一个非常适合IT团队的在线API文档、技术文档工具。搭配的RunApi客户端中一边调试接口、一边自动生成文档。

GoShowDoc 工具通过解析 Go 代码文件中的注释自动生成 RunApi 文档，可通过 RunApi 客户端进行调试。

最棒的是能够解析结构体自动生成参数列表和 JSON 样例。

## 目录

- [命令说明](#命令说明)
- [安装](#安装)
  - [下载可执行文件](#下载可执行文件)
  - [源码编译](#源码编译)
- [设置变量](#设置变量)
- [生成API文档](#生成API文档)
  - [代码示例](#代码示例)
  - [注释格式](#注释格式)
    - [通用API注释](#通用API注释)
    - [API注释](#API注释)
- [生成数据字典](#生成数据字典)
  - [MySQL](#MySQL)
  - [PostgreSQL](#PostgreSQL)
  - [SQLServer](#SQLServer)
  - [SQlite](#SQlite)
- [ShowDoc项目导出导入](#ShowDoc项目导出导入)

## 命令说明

```
$ goshowdoc.exe

NAME:
   goshowdoc - ShowDoc API 接口文档工具

USAGE:
   goshowdoc.exe [global options] command [command options] [arguments...]

VERSION:
   1.2.0

DESCRIPTION:
   项目地址： https://github.com/whaios/goshowdoc
   支持以下功能：
   1. 通过代码注释生成 API 接口文档。
   2. 自动化生成数据字典，支持 mysql, postgres, sqlserver, sqlite3。
   3. 导出和导入 ShowDoc 项目。

COMMANDS:
   flags         查询应用全局相关参数。
   update, u     解析 Go 源码中的注释，生成并更新 ShowDoc 文档。
   datadict, dd  自动生成数据字典。
   item          ShowDoc 项目导出和导入。
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value      ShowDoc 地址。 (default: "https://www.showdoc.com.cn") [%GOSHOWDOC_HOST%]
   --apikey value    ShowDoc 开放 API 认证凭证。 (default: "") [%GOSHOWDOC_APIKEY%]
   --apitoken value  ShowDoc 开放 API 认证凭证。 (default: "") [%GOSHOWDOC_APITOKEN%]
   --debug           开启调试模式。 (default: false)
   --help            显示帮助 (default: false)
   --version, -v     print the version (default: false)
```

## 安装

### 下载可执行文件

已经编译好的平台有： [点击下载](https://github.com/whaios/goshowdoc/releases)
- windows/amd64
- linux/amd64
- darwin/amd64

### 源码编译

如果有 go 开发环境，可以使用如下命令下载 goshowdoc 工具：

```shell
$ go get -u github.com/whaios/goshowdoc

# 1.16 及以上版本
$ go install github.com/whaios/goshowdoc@latest
```

或者下载源码，使用 `go install` 命令构建。

## 设置变量

### 1. ShowDoc 服务器地址

工具中默认配置的官方线上地址 `https://www.showdoc.com.cn` ，如果你也使用的该服务则不需要修改。

如果你使用的是私有版ShowDoc，则使用时需要通过 `--host` 参数指定为自己的地址。

为了避免每次都手动输入该参数，建议将该地址配置为环境变量 `GOSHOWDOC_HOST`。

### 2. 认证凭证

工具生成的文档最终会通过开放API同步到 ShowDoc 的项目中，所以需要配置认证凭证（api_key 和 api_token） 。

登录showdoc > 进入具体项目 > 点击右上角的”项目设置” > “开放API” 便可看到。

为了避免每次生成时都输入这两个参数，建议将该参数配置为环境变量 `GOSHOWDOC_APIKEY` 和 `GOSHOWDOC_APITOKEN`。

## 生成API文档

### 代码示例

```go
package handler

// Handler 书籍管理
//
// Handler 的4个方法，分别对应4个接口文档。
// 这里写的 @catalog @header @resp 三行注释为通用注释，
// 通用注释定义在文件顶部，该文件下的每个接口文档都会包含通用注释。
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
// @query	page		int	true	""	"第几页"
// @query	page_size	int	true	""	"每页显示条数"
// @resp ListRsp{}
func (h *Handler) List() {
}

// Detail 获取指定书籍详情
//
// @url GET {{BASEURL}}/api/v1/book/detail/:id
// @path_var :id int true "" "书籍 id"
// @resp Detail{}
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
// @path_var :id int true "" "书籍 id"
// @remark 危险操作
func (h *Handler) Delete() {
}
```

解析 Go 源码中的注释，生成并更新 ShowDoc 文档。

```shell
# 切换到项目目录
# cd example/ginweb

# 生成并更新文档
$ goshowdoc.exe u --dir ./handler/

# 如果希望输出调试信息，添加 --debug 参数
# goshowdoc.exe --debug u --dir ./handler/
```

输出日志信息：
```
解析Go源码文件 ./example/ginweb/handler/
生成文档(1) 测试文档/书籍/获取书籍列表
生成文档(2) 测试文档/书籍/获取指定书籍详情
生成文档(3) 测试文档/书籍/管理/新建或编辑书籍
生成文档(4) 测试文档/书籍/管理/删除书籍
更新文档 [■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■] 100.00%     [4/4]
更新完成
```

### 注释格式

#### 通用API注释

通用注释定义在文件顶部，该文件下的每个接口文档都会包含通用注释。

注意通用注释的作用范围仅限于 **本文件** 中的接口文档。

| 注释    | 说明    | 示例    |
| ----------------------- | ----------------------- | ----------------------- |
| @catalog | 文档目录，多级目录用 `/` 隔开 | // @catalog 一级/二级/三级 |
| @header | 可选，请求头。格式为 `[字段名] [类型] [必填] ["值"] ["备注"]` | // @header Authorization string true "abc" "用户登录凭证" |
| @response, @resp | 返回内容，支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] ["备注"]`）两种方式。 | // @response TestApiRsp{}  // @param page int "第几页" |
| @response_fail, @resp_fail  | 可选，返回内容。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] ["备注"]`）两种方式。 | // @resp_fail TestApiRsp{}  // @resp_fail page int "第几页" |
| @remark | 可选，备注信息 | // @remark 用户需要先登录 |

#### API注释

| 注释           | 说明         | 示例            |
| ------------- | ----------- | -------------- |
| @title                | 接口文档标题，方法注释。 | // funcName 获取书籍列表 // @title 获取书籍列表  |
| @catalog              | 文档目录，多级目录用 `/` 隔开 | // @catalog 一级/二级/三级 |
| @url                  | 接口URL，格式为：`[method] [url]` | // @url GET {{BASEURL}}/api/v1/book/list |
| @api_status           | 接口状态：0=无，1=开发中，2=测试中，3=已完成，4=需修改，5=已废弃 | // @api_status 3 |
| @description, @desc   | 可选，接口描述信息 | // @description 分页获取书籍列表 |
| @header               | 可选，请求头。格式为 `[字段名] [类型] [必填] ["值"] ["备注"]` | // @header Authorization string true "abc" "用户登录凭证" |
| @path_var             | 可选，请求路径参数。格式为 `[字段名] [类型] [必填] ["值"] ["备注"]` | // @path_var id int true "" "书籍 id" |
| @query                | 可选，请求Query参数。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] [必填] ["值"] ["备注"]`）两种方式。 | // @query id int true "" "书籍 id" |
| @param_mode           | 可选，请求Body参数方式。`urlencoded`、`json` 和 `formdata` | // @param_mode urlencoded |
| @param                | 可选，请求Body参数。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] [必填] ["值"] ["备注"]`）两种方式。 | // @param id int true "" "书籍 id" |
| @response, @resp      | 可选，返回内容。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] ["备注"]`）两种方式。 | // @resp TestApiRsp{}  // @resp page int "第几页" |
| @response_fail, @resp_fail  | 可选，返回内容。支持结构体（如：`Struct{}`，一对大括号结尾） 或 单个参数（如：`[字段名] [类型] ["备注"]`）两种方式。 | // @resp_fail TestApiRsp{}  // @resp_fail page int "第几页" |
| @remark               | 可选，备注信息 | // @remark 用户需要先登录 |

## 生成数据字典

```
$ .\goshowdoc-windows-amd64.exe dd --help

NAME:
   goshowdoc-windows-amd64.exe datadict - 自动生成数据字典。

USAGE:
   goshowdoc-windows-amd64.exe datadict [command options] [arguments...]

OPTIONS:
   --driver value          数据库类型，支持：mysql, postgres, sqlserver, sqlite3 (default: "mysql")
   --host value, -h value  数据库地址和端口，如果是SQlite数据库则为文件 (default: "127.0.0.1:3306")
   --user value, -u value  数据库用户名
   --pwd value, -p value   数据库密码
   --db value              要同步的数据库名
   --schema value          PostgreSQL 数据库模式 (default: "public")
   --cat value             文档所在目录，如果需要多层目录请用斜杠隔开，例如：“一层/二层/三层”
   --help                  显示帮助 (default: false)
```

### MySQL

```shell
.\goshowdoc-windows-amd64.exe dd --driver mysql -h 127.0.0.1:3306 -u root -p 123456 --db testdb
```

### PostgreSQL

```shell
.\goshowdoc-windows-amd64.exe dd --driver postgres -h 127.0.0.1:5432 -u postgres -p 123456 --db postgres
```

### SQLServer

```shell
.\goshowdoc-windows-amd64.exe dd --driver sqlserver -h 127.0.0.1:1433 -u sa -p 123456 --db testdb
```

### SQlite

因为 go-sqlite3 库是一个 cgo 库，需要 gcc 环境，编译的可执行程序默认没有开启 CGO，所以无法正常连接 SQlite 数据库。
如果有需要，请自行下载代码编译。

```shell
.\goshowdoc-windows-amd64.exe dd --driver sqlite3 -h .\test.db
```

## ShowDoc项目导出导入

项目导入导出功能使用 ShowDoc 内部接口而非开放接口，所以可能会出现ShowDoc内部接口调整导致导入功能无法正常使用。

目前兼容 ShowDoc 版本为 **v2.9.14**