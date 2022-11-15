package initialize

import (
	"cloud_disk/greet/define"
	"cloud_disk/greet/internal/config"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql" ///对应数据库的驱动必须要打上
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"net/url"
	"os"
	"xorm.io/xorm"
)

// Init 初始化mysql数据连接对象
func Init(dataSource string) *xorm.Engine {
	//driverName := "mysql"
	//dataSourceName := "用户名:密码@/数据库名称?charset=utf8"
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log.Printf("获取mysql连接对象失败：%v", err)
		return nil
	}
	return engine
}

// InitRedis 获取redis的连接对象
func InitRedis(c config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
	})
}

// InitTenCosClient 初始化获取TenCos连接
func InitTenCosClient() *cos.Client {
	u, _ := url.Parse(define.TenCosURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥和ID
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv(define.TenSecretID), //我是放在环境变量中获取的（自己可以去自己的腾讯云服务获取，上面注释有提示）
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv(define.TenSecretKey),
		},
	})
	return client
}
