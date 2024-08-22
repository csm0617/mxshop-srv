package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop_srvs/user_srv/global"
)

/**
 * 获取环境变量
 * @param env 环境变量名
 * @return bool 环境变量值，true为true, false为false
 */
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	//从环境变量中拿到配置文件应该开启哪个
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePath := "user_srv/"
	configFilePrefix := "config"
	//默认开启生产环境
	configFileName := fmt.Sprintf("%s%s-pro.yaml", configFilePath, configFilePrefix)
	if debug { //否则开发环境
		configFileName = fmt.Sprintf("%s%s-debug.yaml", configFilePath, configFilePrefix)
	}
	v := viper.New()                //拿到一个viper的实例
	v.SetConfigFile(configFileName) // 设置配置文件路径
	if err := v.ReadInConfig(); err != nil {
		zap.S().Error("read config file failed!")
		panic(err)
	}
	//将配置文件中的配置读到全局变量中
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		zap.S().Error("unmarshal SeverConfig failed!")
		panic(err)
	}
	zap.S().Infof("succeed unmarshaling ServerConfig:%v", global.ServerConfig)

	//!!!动态监听配置文件的变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) { //（配置文件发生改变）触发通知事件，需要做什么事
		zap.S().Infof("Config file changed: %s", e.Name)
		//1.重新读取配置文件
		_ = v.ReadInConfig()
		//2.重新将配置文件解析到全局的配置文件中
		zap.S().Infof("reloading config file: %s", e.Name)
		if err := v.Unmarshal(global.ServerConfig); err != nil {
			zap.S().Error("unmarshal SeverConfig failed!")
			panic(err)
		}
		zap.S().Infof("succeed unmarshaling ServerConfig:%v", global.ServerConfig)
		//todo: 重启应用或其他需要的操作
		//...
	})
}
