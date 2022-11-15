package models

import "time"

type UserRepository struct {
	Id                 int
	Identity           string
	UserIdentity       string
	ParentId           int64 //记录父层级ID
	RepositoryIdentity string
	Ext                string
	Name               string
	CreatedAt          time.Time `xorm:"created_at"`
	UpdatedAt          time.Time `xorm:"updated_at"`
	DeletedAt          time.Time `xorm:"deleted_at"`
}

func (table UserRepository) TableName() string {
	return "user_repository"
}
