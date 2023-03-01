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
	CreationTime   Time `gorm:"creation_time;type:timestamp;default:current_timestamp" json:"creation_time"`
	LastModifyTime Time `gorm:"last_modify_time;type:timestamp" json:"last_modify_time"`
}
