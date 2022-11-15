package models

import "time"

// ShareBasic 用户文件分享基本信息
type ShareBasic struct {
	Id                     int
	Identity               string
	UserIdentity           string
	UserRepositoryIdentity string
	RepositoryIdentity     string
	ExpiredTime            int       //分享过期时间
	ClickNum               int       //分享点击数
	CreatedAt              time.Time `xorm:"created_at"`
	UpdatedAt              time.Time `xorm:"updated_at"`
	DeletedAt              time.Time `xorm:"deleted_at"`
}

func (table ShareBasic) TableName() string {
	return "share_basic"
}
