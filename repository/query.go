package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sailing.cn/v2/pager"
)

// IQuery 查询接口
type IQuery interface {
	Query() map[string]interface{}
}

// IPageQuery 分页查询接口
type IPageQuery interface {
	Query() map[string]interface{}
	GetPageQuery() *pager.PageQuery
}

// Query 查询
func (c *DbContext) Query(result interface{}, order string, query IQuery) error {
	_, err := c.query(result, order, 0, 0, query, false)
	return err
}

// PageQuery 分页查询
func (c *DbContext) PageQuery(query IPageQuery) (result *pager.PageResult, err error) {
	pagerQuery := query.GetPageQuery()
	var total int64 = 0
	list := make([]interface{}, 0)
	if result == nil {
		result = &pager.PageResult{}
	}
	offset := pagerQuery.PageSize * (pagerQuery.Page - 1) //mysql处理逻辑
	total, err = c.query(&list, pagerQuery.Order, int(pagerQuery.PageSize), int(offset), query, true)
	if err != nil {
		return nil, err
	}
	result.Build(list, total, pagerQuery)
	return result, nil
}

func (c *DbContext) query(result interface{}, order string, limit, offset int, query IQuery, page bool) (total int64, err error) {
	//todo 将所有的分页查询设置为按创建时间排序
	if len(order) == 0 {
		order = "creation_time desc"
	}
	var (
		table Table
	)
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr {
		return 0, errors.New("result is not a pointer")
	}

	sliceType := resultType.Elem()
	if sliceType.Kind() != reflect.Slice {
		return 0, errors.New("result doesn't point to a slice")
	}
	// *Item
	itemPtrType := sliceType.Elem()
	// Item
	itemType := itemPtrType.Elem()

	elemValue := reflect.New(itemType)
	elemValueType := reflect.TypeOf(elemValue)
	tableType := reflect.TypeOf((*Table)(nil)).Elem()

	if elemValueType.Implements(tableType) {
		return 0, errors.New("neither the query nor result implement Table")
	}

	table = elemValue.Interface().(Table)
	db := c.Table(table.TableName())
	if query != nil {
		m := query.Query()
		for item := range m {
			db = db.Where(item, m[item])
		}
	}

	if len(order) != 0 {
		db = db.Order(order)
	}
	if page {
		db.Count(&total)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	err = db.Scan(result).Error
	if err != nil {
		log.Errorf("查询语句执行失败:%s", err.Error())
	}
	return
}
