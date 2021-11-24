package main

import (
	"fmt"

	envParser "github.com/caarlos0/env"
)

type Settings struct {
	ServiceProtocol string
	ServicePort     string

	LogFilePath string

	RedisProtocol string
	RedisPort     string
}

type EnvironmentVariables struct {
	REDIS_PORT           string `env:"REDIS_PORT,required" envDefault:":6380"`
	SESSION_LOG_LOCATION string `env:"SESSION_LOG_LOCATION" envDefault:"/var/log/emailLogs.log"`
}

func InitSettings() (settings Settings) {
	env := EnvironmentVariables{}
	if err := envParser.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	settings = Settings{
		ServiceProtocol: "tcp",
		ServicePort:     "0.0.0.0:8082",

		LogFilePath: env.SESSION_LOG_LOCATION,

		RedisProtocol: "tcp",
		RedisPort:     fmt.Sprintf("redis:%s", env.REDIS_PORT),
	}

	return
}
