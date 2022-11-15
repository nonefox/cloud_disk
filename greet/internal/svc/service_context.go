package svc

import (
	"cloud_disk/greet/initialize"
	"cloud_disk/greet/internal/config"
	"cloud_disk/greet/internal/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"xorm.io/xorm"
)

// ServiceContext 服务需要的一些对象放到这里集中配置
type ServiceContext struct {
	Config      config.Config
	MysqlEngine *xorm.Engine
	RedDB       *redis.Client
	Auth        rest.Middleware
}

// NewServiceContext 把读出的yaml的配置信息拿到，实例化
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		MysqlEngine: initialize.Init(c.Mysql.DataSource),
		RedDB:       initialize.InitRedis(c),
		Auth:        middleware.NewAuthMiddleware().Handle,
	}
}
