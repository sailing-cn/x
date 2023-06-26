package mysql

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sailing.cn/v2/repository/conf"
)

func NewMysql(cnf *conf.Config) *gorm.DB {
	if len(cnf.Db.Connection) == 0 {
		panic("数据库连接字符未配置")
	}
	context, err := gorm.Open(mysql.Open(cnf.Db.Connection), &gorm.Config{})
	if err != nil {
		log.Errorf("数据库打开失败 %v", err.Error())
	}
	return context
}
