package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	"mxshop_srvs/user_srv/global"
)

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 2, // Slow SQL threshold
			LogLevel:      logger.Info,     // Log level
			Colorful:      true,            // Disable color
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.ServerConfig.MysqlInfo.User,
		global.ServerConfig.MysqlInfo.Password,
		global.ServerConfig.MysqlInfo.Host,
		global.ServerConfig.MysqlInfo.Port)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, //迁移时禁用创建外键约束
	})
	if err != nil {
		panic(err)
	}

}
