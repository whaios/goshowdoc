package review

import "ginweb/model/book"

// Review 书籍评论
type Review struct {
	Id           book.Id `json:"id"`            // 评论id
	CreationUnix int64   `json:"creation_unix"` // 发表时间
	BookId       book.Id `json:"book_id"`       // 书籍id
	Content      string  `json:"content"`       // 评论内容
	Reviewer             // 评论人

	RecursiveReviews []*Review `json:"recursive_reviews"` // 测试是否能安全解析递归类型
}

// Reviewer 评论人
type Reviewer struct {
	ReviewUserId   int64  `json:"review_user_id"`   // 评论人id
	ReviewUserName string `json:"review_user_name"` // 评论人名称
}
