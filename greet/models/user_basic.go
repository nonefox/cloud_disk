package models

import "time"

// UserBasic 网盘用户数据基本结构
type UserBasic struct {
	Id        int
	Identity  string
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time `xorm:"created_at"`
	UpdatedAt time.Time `xorm:"updated_at"`
	DeletedAt time.Time `xorm:"deleted_at"`
}

func (table UserBasic) TableName() string {
	return "user_basic"
}
