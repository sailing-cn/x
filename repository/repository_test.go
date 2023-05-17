package repository

import (
	"context"
	"sailing.cn/pager"
	"testing"
)

func TestDm(t *testing.T) {
	//dsn := "dm://SAILING:Sl123.com@192.168.1.65:5236?autoCommit=true"
	//InitWithCnf(&conf_d.Config{Db: &conf_d.DB{
	//	Connection: dsn,
	//	Type:       "dm",
	//	Debug:      true,
	//}})
	//GetContext().AutoMigrate(&Client{})

	var r interface{}

	q := TestPageQuery{}
	_, err := GetContext().Context(context.Background()).PageListQuery(&r, &q)
	if err != nil {
		return
	}

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
	Id           string `json:"id,omitempty" gorm:"index:key,unique;type:varchar(255);primaryKey;comment:唯一标识"` //
	Name         string `json:"name,omitempty" gorm:"type:varchar(255)"`                                        //
	Icon         string `json:"icon,omitempty" gorm:"type:varchar(255)"`                                        //
	Isv          string `json:"isv,omitempty" gorm:"type:varchar(255)"`                                         //
	CreationTime int64  `json:"creation_time,omitempty" gorm:"type:bigint"`                                     //
}

type TestPageQuery struct {
	Client
	pager.PageQuery
}

func (t *TestPageQuery) Query() map[string]interface{} {
	return nil
}
