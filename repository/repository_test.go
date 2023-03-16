package repository

import (
	"gorm.io/gorm"
	conf_d "sailing.cn/repository/conf.d"
	dmSchema "sailing.cn/repository/driver/dm/schema"
	"testing"
	"time"
)

func TestDm(t *testing.T) {
	dsn := "dm://SAILING:Sl123.com@192.168.1.65:5236?autoCommit=true"
	InitWithCnf(&conf_d.Config{Db: &conf_d.DB{
		Connection: dsn,
		Type:       "dm",
		Debug:      true,
	}})
	GetContext().AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Key      string `gorm:"index:key,unique"`
	Name     string `gorm:"index:name"`
	Age      int
	Content  dmSchema.Clob `gorm:"size:1024000"`
	Birthday time.Time
}
