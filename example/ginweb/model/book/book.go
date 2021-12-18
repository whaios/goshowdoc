package book

// Book 书籍的编辑对象
type Book struct {
	Id        int64  `json:"id,string"`                 // id
	Title     string `json:"title" validate:"required"` // 书名
	Type      string `json:"type"`                      // 包装：平装、精装
	Pages     int    `json:"pages" validate:"min=1"`    // 页数
	PubDate   int64  `json:"pub_date"`                  // 出版日期
	Publisher string `json:"publisher"`                 // 出版社
	Isbn      string `json:"isbn"`                      // 图书编号
	IsActive  bool   `json:"is_active"`                 // 是否激活
}
