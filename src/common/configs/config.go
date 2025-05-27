package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode string `mapstructure:"MODE"`

	Server struct {
		Name    string `mapstructure:"SERVER_NAME"`
		Address string `mapstructure:"SERVER_ADDRESS"`
		Prefix  string `mapstructure:"SERVER_PREFIX"`
	} `mapstructure:",squash"`

	Jwt struct {
		AccessSecret  string `mapstructure:"JWT_ACCESS_SECRET"`
		RefreshSecret string `mapstructure:"JWT_REFRESH_SECRET"`
		ExpireAccess  int64  `mapstructure:"JWT_EXPIRE_ACCESS"`
		ExpireRefresh int64  `mapstructure:"JWT_EXPIRE_REFRESH"`
	} `mapstructure:",squash"`

	Postgres struct {
		Host        string `mapstructure:"POSTGRES_HOST"`
		Port        int    `mapstructure:"POSTGRES_PORT"`
		User        string `mapstructure:"POSTGRES_USER"`
		Password    string `mapstructure:"POSTGRES_PASSWORD"`
		DB          string `mapstructure:"POSTGRES_DB"`
		SslMode     string `mapstructure:"POSTGRES_SSL_MODE"`
		AutoMigrate bool   `mapstructure:"POSTGRES_AUTO_MIGRATE"`
		MaxLifeTime int    `mapstructure:"POSTGRES_MAX_LIFE_TIME"`
	} `mapstructure:",squash"`

	Redis struct {
		Address  string `mapstructure:"REDIS_ADDRESS"`
		Password string `mapstructure:"REDIS_PASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
	} `mapstructure:",squash"`

	Mail struct {
		ClientID     string `mapstructure:"MAIL_CLIENT_ID"`
		ClientSecret string `mapstructure:"MAIL_CLIENT_SECRET"`
		RefreshToken string `mapstructure:"MAIL_REFRESH_TOKEN"`
		AccessToken  string `mapstructure:"MAIL_ACCESS_TOKEN"`
	} `mapstructure:",squash"`

	Minio struct {
		User      string `mapstructure:"MINIO_USER"`
		Password  string `mapstructure:"MINIO_PASSWORD"`
		AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
		SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
		Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
		Bucket    string `mapstructure:"MINIO_BUCKET"`
	} `mapstructure:",squash"`
}

var common *Config

func GetConfig() *Config {
	return common
}

func InitConfig(path string) error {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&common)
	return err
}
