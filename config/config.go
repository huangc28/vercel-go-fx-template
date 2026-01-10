package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName   string `mapstructure:"app_name"`
	AppEnv    string `mapstructure:"app_env"`
	AppPort   int    `mapstructure:"app_port"`
	LogLevel  string `mapstructure:"log_level"`
	RedisURL  string `mapstructure:"redis_url"`
	PGURL     string `mapstructure:"pg_url"`
	InngestID string `mapstructure:"inngest_app_id"`
}

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("app_name", "vercel-go-service")
	v.SetDefault("app_env", "development")
	v.SetDefault("app_port", 3010)
	v.SetDefault("log_level", "info")

	v.SetDefault("redis_url", "")
	v.SetDefault("pg_url", "")

	v.SetDefault("inngest_app_id", "vercel-go-service")
	return v
}

func NewConfig(v *viper.Viper) (Config, error) {
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
