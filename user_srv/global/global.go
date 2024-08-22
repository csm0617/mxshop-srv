package global

import (
	"gorm.io/gorm"

	"mxshop_srvs/user_srv/config"
)

var (
	DB           *gorm.DB             //定义全局变量db
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
