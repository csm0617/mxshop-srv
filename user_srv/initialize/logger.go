package initialize

import (
	"go.uber.org/zap"
)

func InitLogger() {
	//todo:设置日志配置等
	//设置日志级别为生产环境
	logger, _ := zap.NewDevelopment()
	//替换zap全局的Logger
	zap.ReplaceGlobals(logger)
}
