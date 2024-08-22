package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	Host      string      `mapstructure:"host" json:"host"`
	Port      int         `mapstructure:"port" json:"port"`
	MysqlInfo MysqlConfig `mapstructure:"mysql" json:"mysql"`
}
