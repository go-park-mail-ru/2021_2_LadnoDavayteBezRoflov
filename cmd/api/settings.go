package main

import (
	"fmt"
	"time"

	envParser "github.com/caarlos0/env"
	"github.com/gin-contrib/cors"
)

type Settings struct {
	RootURL    string
	SessionURL string
	ProfileURL string
	BoardsURL  string

	ServerAddress string

	Origins        []string
	AllowedMethods []string

	SessionCookieLifeTimeInDays time.Duration

	corsConfig cors.Config

	LogFilePath string

	RedisProtocol string
	RedisPort     string

	PostgresDsn string
}

type EnvironmentVariables struct {
	DB_PORT           string `env:"DB_PORT",required`
	REDIS_PORT        string `env:"REDIS_PORT",required`
	POSTGRES_USER     string `env:"POSTGRES_USER",required`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD",required`
	DATABASE_HOST     string `env:"DATABASE_HOST,required"`
	POSTGRES_DB       string `env:"POSTGRES_DB,required"`
	FRONTEND_ADDRESS  string `env:"FRONTEND_ADDRESS,required"`
	LOG_LOCATION      string `env:"LOG_LOCATION" envDefault:"/var/log/backendLogs.log"`
}

func InitSettings() (settings Settings) {
	env := EnvironmentVariables{}
	if err := envParser.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	settings = Settings{
		RootURL:    "/api",
		SessionURL: "/sessions",
		ProfileURL: "/profile",
		BoardsURL:  "/boards",

		ServerAddress: ":8000",

		Origins: []string{
			"http://localhost:8000",
			fmt.Sprintf("http://%s", env.FRONTEND_ADDRESS),
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		SessionCookieLifeTimeInDays: 3,

		corsConfig: cors.DefaultConfig(),

		LogFilePath: env.LOG_LOCATION,

		RedisProtocol: "tcp",
		RedisPort:     fmt.Sprintf("redis:%s", env.REDIS_PORT),

		PostgresDsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.DATABASE_HOST, env.POSTGRES_USER, env.POSTGRES_PASSWORD, env.POSTGRES_DB, env.DB_PORT),
	}

	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	return
}
