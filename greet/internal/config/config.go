package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	//把数据库的配置信息全部抽取到一起，从yaml配置文件中读取
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr string
	}
}
