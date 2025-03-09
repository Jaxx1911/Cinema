package configs

import "github.com/spf13/viper"

type Config struct {
	Mode string `mapstructure:"mode"`

	Server struct {
		Name    string `mapstructure:"name"`
		Address string `mapstructure:"address"`
		Prefix  string `mapstructure:"prefix"`
	} `mapstructure:"server"`
	Postgres struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		User        string `mapstructure:"user"`
		Password    string `mapstructure:"password"`
		DB          string `mapstructure:"db"`
		SslMode     string `mapstructure:"ssl_mode"`
		AutoMigrate bool   `mapstructure:"auto_migrate"`
		MaxLifeTime int    `mapstructure:"max_life_time"`
	}
	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

var common *Config

func GetConfig() *Config {
	return common
}

func InitConfig(path string) error {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&common)

	return err
}
