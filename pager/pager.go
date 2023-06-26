package pager

// Page 分页
type Page struct {
	Total int64 `json:"total"` //总数
	Page  int32 `json:"page"`  //页索引
	Count int32 `json:"-"`
	Size  int32 `json:"size"` //分页大小
}

// PageQuery 分页查询
type PageQuery struct {
	Page     int32  `json:"page" form:"page"`           //页索引
	PageSize int32  `json:"page_size" form:"page_size"` //分页大小
	Order    string `json:"order" form:"order"`         //排序
}

// PageResult 分页结果
type PageResult struct {
	Data interface{} `json:"data"` //数据列表
	Page *Page       `json:"page"` //分页信息
}
