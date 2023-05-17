package repository

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sailing.cn/pager"
	"sailing.cn/repository/ch"
	conf "sailing.cn/repository/conf.d"
	"sailing.cn/repository/dm"
	"sailing.cn/repository/mysql"
	"sailing.cn/repository/postgres"
	"sailing.cn/repository/sqlite"
	"sailing.cn/utils"
)

type IQuery interface {
	Query() map[string]interface{}
}

type IPageQuery interface {
	Query() map[string]interface{}
	GetPageQuery() *pager.PageQuery
}

var (
	_context = &Context{}
	cnf      = &conf.Config{}
)

type Context struct {
	*gorm.DB
}

type Table interface {
	TableName() string
}

func new() {
	switch cnf.Db.Type {
	case "mysql":
		_context.DB = mysql.NewMysql(cnf)
		break
	case "clickhouse":
		_context.DB = ch.NewClickHouse(cnf)
		break
	case "sqlite":
		_context.DB = sqlite.NewSqlite(cnf)
	case "dm":
		_context.DB = dm.NewDm(cnf)
	case "postgres", "pgsql":
		_context.DB = postgres.NewPostgres(cnf)
	default:
		log.Errorf("暂不支持该数据类型:%s", cnf.Db.Type)
		break
	}
	if cnf.Db.Debug {
		_context.DB = _context.Debug()
	}
	_context.DB.Use(tracing.NewPlugin())
}

func GetContext() *Context {
	if _context == nil || _context.DB == nil {
		new()
	}
	return _context
}

func (c *Context) Context(ctx context.Context) *Context {
	return &Context{c.DB.WithContext(ctx)}
}

func InitWithCnf(_cnf *conf.Config) {
	cnf = _cnf
}

func Init() {
	path := filepath.Join(utils.GetExecPath(), "conf.d", "conf.yml")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, cnf)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}

func (c *Context) TableX(table string) *Context {
	c.Table(table)
	return c
}

func (c *Context) IsExist(table string, query interface{}, args ...interface{}) (bool, error) {
	var count int64 = 0
	result := GetContext().Table(table).Where(query, args...).Count(&count)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	}
	return count > 0, result.Error
}

func (c *Context) PageListQuery(result interface{}, query IPageQuery) (page *pager.PageList, err error) {
	pagerQuery := query.GetPageQuery()
	var total int64 = 0
	if page == nil {
		page = &pager.PageList{}
	}
	offset := pagerQuery.PageSize * (pagerQuery.Page - 1) //mysql处理逻辑
	total, err = c.list2(result, pagerQuery.Order, int(pagerQuery.PageSize), int(offset), query, true)
	if err != nil {
		return nil, err
	}
	page.ToPager(result, total, pagerQuery)
	return page, nil
}

// PageList deprecated
// 建议使用PageListQuery
func (c *Context) PageList(result interface{}, pagerQuery *pager.PageQuery, query map[string]interface{}) (page *pager.PageList, err error) {
	pagerQuery.GetQuery()
	var total int64 = 0
	if page == nil {
		page = &pager.PageList{}
	}
	total, err = c.list(result, pagerQuery.Order, int(pagerQuery.PageSize), int(pagerQuery.Page), query, true)
	if err != nil {
		return nil, err
	}
	page.ToPager(result, total, pagerQuery)
	return page, nil
}

func (c *Context) ListQuery(result interface{}, order string, query IQuery) (err error) {
	_, err = c.list2(result, order, 0, 0, query, false)
	return err
}

// List deprecated
// 建议使用ListQuery
func (c *Context) List(result interface{}, order string, query map[string]interface{}) (err error) {
	_, err = c.list(result, order, 0, 0, query, false)
	return err
}

func (c *Context) Preload(query string, args ...interface{}) *Context {
	return &Context{c.DB.Preload(query, args)}
}
func (c *Context) Model(value interface{}) *Context {
	return &Context{c.DB.Model(value)}
}

// deprecated
func (c *Context) list(result interface{}, order string, limit, offset int, query interface{}, page bool) (total int64, err error) {
	//todo 将所有的分页查询设置为按创建时间排序
	if len(order) == 0 {
		order = "creation_time desc"
	}
	var (
		table Table
		ok    bool
	)
	if table, ok = query.(Table); !ok {
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
	}
	db := c.Table(table.TableName())
	if query != nil {
		db = db.Where(query)
	}

	if len(order) != 0 {
		db = db.Order(order)
	}
	if offset > 0 {
		db = db.Offset(offset - 1)
	}
	if page {
		db.Count(&total)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}

	err = db.Find(result).Error
	if err != nil {
		log.Errorf("查询语句执行失败:%s", err.Error())
	}
	return
}

func (c *Context) list2(result interface{}, order string, limit, offset int, query IQuery, page bool) (total int64, err error) {
	//todo 将所有的分页查询设置为按创建时间排序
	if len(order) == 0 {
		order = "creation_time desc"
	}
	var (
		table Table
		//ok    bool
	)
	//if table, ok = query.(Table); !ok {
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
	//}
	db := c.Table(table.TableName())
	if query != nil {
		m := query.Query()
		for item, _ := range m {
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
