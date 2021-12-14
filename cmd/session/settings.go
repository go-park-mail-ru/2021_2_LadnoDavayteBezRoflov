package main

import (
	"fmt"

	envParser "github.com/caarlos0/env"
	"github.com/spf13/viper"
)

type Settings struct {
	ServiceProtocol string `mapstructure:"service_protocol"`
	ServicePort     string `mapstructure:"service_port"`

	ServiceMetricsPort string `mapstructure:"service_metrics_port"`

	LogFilePath string

	RedisProtocol string `mapstructure:"redis_protocol"`
	RedisPort     string
}

type EnvironmentVariables struct {
	REDIS_PORT           string `env:"REDIS_PORT,required" envDefault:":6380"`
	SESSION_LOG_LOCATION string `env:"SESSION_LOG_LOCATION" envDefault:"/var/log/sessionServiceLogs.log"`
}

func InitSettings() (settings Settings) {
	env := EnvironmentVariables{}
	if err := envParser.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	viper.AddConfigPath("./cmd/session")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%+v\n", err)
	}

	settings = Settings{
		ServiceProtocol: viper.GetString("service_protocol"),
		ServicePort:     viper.GetString("service_port"),

		ServiceMetricsPort: viper.GetString("service_metrics_port"),

		LogFilePath: env.SESSION_LOG_LOCATION,

		RedisProtocol: viper.GetString("redis_protocol"),
		RedisPort:     fmt.Sprintf("redis:%s", env.REDIS_PORT),
	}

	return
}
