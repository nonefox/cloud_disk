package define

import (
	"cloud_disk/greet/initialize"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

var MailPassword = os.Getenv("MailPassword")

//// MysqlEngine 定义一个初始化mysql连接变量(这里就不再需要数据库配置了，在srv中统一配置)
//var MysqlEngine = initialize.Init(config.Config{}.Mysql.DataSource)
//
//// RedDB 定义一个redis初始化连接变量
//var RedDB = initialize.InitRedis(config.Config{})

// JWTKey 定义JWTToken使用的key
var JWTKey = []byte("cloud_disk")

// UserClaim 用户的声明结构（我们会对他进行签名生成用户的token信息）
type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.RegisteredClaims
}

// 腾讯云对象存储
var TenSecretID = os.Getenv("TencentSecretID")
var TenSecretKey = os.Getenv("TencentSecretKey")
var TenCosURL = "https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com" //使用自己的
var TenCosClient = initialize.InitTenCosClient()

// PageSize 分页的默认参数
var PageSize = 10
var PageIndex = 1

var Datetime = "2008-08-08 08:08:08"

//token的有效时间
var TokenTime = 3600

//刷新token时间
var RefreshTokenTime = 7200
