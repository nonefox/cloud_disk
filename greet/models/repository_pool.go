package models

import "time"

// RepositoryPool 用户文件存储中心
type RepositoryPool struct {
	Id        int
	Identity  string
	Hash      string //文件的hash值
	Name      string
	Ext       string
	Size      int64
	Path      string
	CreatedAt time.Time `xorm:"created_at"`
	UpdatedAt time.Time `xorm:"updated_at"`
	DeletedAt time.Time `xorm:"deleted_at"`
}

func (table RepositoryPool) TableName() string {
	return "repository_pool"
}
