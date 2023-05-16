package pager

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
)

type BasePageQuery interface {
	GetPageQuery() *PageQuery
	ToPage() *PageList
}

type Page struct {
	Total int64 `json:"total"` //总数
	Page  int32 `json:"page"`  //页索引
	Count int32 `json:"-"`
	Size  int32 `json:"size"` //分页大小
}

type PageQuery struct {
	Page     int32  `json:"page" form:"page"`            //页索引
	PageSize int32  `json:"page_size" form:"page_size""` //分页大小
	Order    string `json:"order" form:"order"`          //排序
}

type PageList struct {
	Data interface{} `json:"data"` //数据列表
	Page *Page       `json:"page"` //分页信息
}

func NewPageQuery(page int32, size int32, order string) *PageQuery {
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	if len(order) <= 0 {
		order = "creation_time desc"
	}
	return &PageQuery{Page: page, PageSize: size, Order: order}
}

func GetPageQuery(page *wrapperspb.Int32Value, size *wrapperspb.Int32Value, order *wrapperspb.StringValue) *PageQuery {
	result := &PageQuery{}
	if size != nil {
		result.PageSize = (size.GetValue())
	} else {
		result.PageSize = 20
	}
	if page != nil {
		result.Page = (page.GetValue())
	} else {
		result.Page = 1
	}
	if order != nil {
		result.Order = order.Value
	}
	return result
}

func (p *PageQuery) GetQuery() {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.Page <= 0 {
		p.Page = 1
	}
}

func (p *PageList) ToPager(source interface{}, total int64, query *PageQuery) {
	page := &Page{
		Total: total,
		Page:  int32(query.Page),
		Count: 0,
		Size:  int32(query.PageSize),
	}
	p.Data = source
	p.Page = page
	if page.Total == 0 || page.Size == 0 {
		page.Count = 0
	} else {
		if float64(page.Total)/float64(page.Size) == 0 {
			page.Count = int32(math.Ceil(float64(page.Total) / float64(page.Size)))
		} else {
			page.Count = int32(math.Ceil(float64(page.Total)/float64(page.Size))) + 1
		}
	}
}
