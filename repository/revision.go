package repository

import (
	"database/sql/driver"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type Revision []byte

//func (r Revision) GetVersion() string {
//	return r.Revision
//}
//
//func (r Revision) SetVersion() {
//	r.Revision = uuid.NewV4().String()
//}

func (r *Revision) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	*r = bytes
	return nil
}

// Value return json value, implement driver.Valuer interface
func (r *Revision) Value() (driver.Value, error) {
	return []byte(uuid.NewV4().String()), nil
}
