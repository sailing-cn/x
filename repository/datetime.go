package repository

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = Time(t1)
	return err
}
func (t *Time) MarshalJSON() ([]byte, error) {
	// tune := fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))
	tune := time.Time(*t).Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t *Time) Value() (driver.Value, error) {
	return time.Time(*t), nil
}

func (t *Time) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = Time(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

type BaseModel struct {
	CreationTime int64     `gorm:"column:creation_time;type:bigint;not null;comment:创建时间"` //创建时间
	Revision     *Revision `gorm:"column:revision;type:varbinary;not null;comment:乐观锁"`    //乐观锁
	IsDeleted    bool      `gorm:"column:is_deleted;type:bit;not null;comment:逻辑删除"`       //逻辑删除
}
