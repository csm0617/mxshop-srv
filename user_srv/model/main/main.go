package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"io"
	"strings"
)

func genMd5(code string) string {
	Md5 := md5.New()                 //创建md5对象
	_, _ = io.WriteString(Md5, code) //将字符串写入Md5进行加密
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second * 2, // Slow SQL threshold
	//		LogLevel:      logger.Info,     // Log level
	//		Colorful:      true,            // Disable color
	//	},
	//)
	//dsn := "root:123456@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: newLogger,
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true, //使用单数表名
	//	},
	//	DisableForeignKeyConstraintWhenMigrating: true, //迁移时禁用创建外键约束
	//})
	//if err != nil {
	//	panic(err)
	//}
	//_ = db.AutoMigrate(&model.User{})
	//fmt.Println(genMd5("a1234561"))
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodePwd := password.Encode("123456", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodePwd)
	passwordInfo := strings.Split(newPassword, "$")
	fmt.Println(passwordInfo)
	verify := password.Verify("123456", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(verify)
}
