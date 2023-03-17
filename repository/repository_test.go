package repository

import (
	conf_d "sailing.cn/repository/conf.d"
	"testing"
)

func TestDm(t *testing.T) {
	dsn := "dm://SAILING:Sl123.com@192.168.1.65:5236?autoCommit=true"
	InitWithCnf(&conf_d.Config{Db: &conf_d.DB{
		Connection: dsn,
		Type:       "dm",
		Debug:      true,
	}})
	GetContext().AutoMigrate(&Client{})
}

//type User struct {
//	gorm.Model
//	Key      string `gorm:"index:key,unique"`
//	Name     string `gorm:"index:name"`
//	Age      int
//	Content  dmSchema.Clob `gorm:"size:1024000"`
//	Birthday time.Time
//}

type Client struct {
	Id           string `json:"id,omitempty" gorm:"index:key,unique;type:varchar(255);primaryKey;"` //
	Name         string `json:"name,omitempty" gorm:"type:varchar(255)"`                            //
	Icon         string `json:"icon,omitempty" gorm:"type:varchar(255)"`                            //
	Isv          string `json:"isv,omitempty" gorm:"type:varchar(255)"`                             //
	CreationTime int64  `json:"creation_time,omitempty" gorm:"type:bigint"`                         //
}
