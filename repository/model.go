package repository

type BaseModel struct {
	CreationTime int64     `gorm:"column:creation_time;type:bigint;not null;comment:创建时间"` //创建时间
	Revision     *Revision `gorm:"column:revision;type:varbinary;not null;comment:乐观锁"`    //乐观锁
	IsDeleted    bool      `gorm:"column:is_deleted;type:bit;not null;comment:逻辑删除"`       //逻辑删除
}
