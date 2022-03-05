package book

import (
	. "ginweb/comm"               // 测试 . 包
	"ginweb/model/book"           // 测试正常导入包
	review1 "ginweb/model/review" // 测试包别名
)

// ListRsp 列表返回结果
type ListRsp struct {
	TotalCount int         `json:"total_count"` // 总条数
	Items      []*ListItem `json:"items"`       // 书籍
}

// ListItem 列表项
type ListItem struct {
	Id        book.Id  `json:"id,string"` // 标识符
	Title     string   `json:"title"`     // 书名
	Publisher string   `json:"publisher"` // 出版社
	Tags      []string `json:"tags"`      // 标签
}

// Detail 书籍详情返回结果
type Detail struct {
	book.Book
	PubDateStr string            `json:"pub_date_str"` // 出版日期
	Reviews    []*review1.Review `json:"reviews"`      // 书籍评论
	ReviewPage *Page             `json:"review_page"`  // 书籍评论分页
}
