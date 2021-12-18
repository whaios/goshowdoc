package book

import (
	_ "ginweb/comm"
	_ "ginweb/model/book"
)

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
