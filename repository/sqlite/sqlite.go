package sqlite

import (
	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	conf "sailing.cn/repository/conf.d"
)

func NewSqlite(cnf *conf.Config) *gorm.DB {
	if len(cnf.Db.Connection) == 0 {
		panic("数据库连接字符未配置")
	}
	context, err := gorm.Open(sqlite.Open(cnf.Db.Connection), &gorm.Config{})
	if err != nil {
		log.Errorf("数据库打开失败 %v", err.Error())
	}
	return context
}
