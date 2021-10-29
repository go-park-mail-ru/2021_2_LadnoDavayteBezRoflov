package main

import (
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

	corsConfig cors.Config

	LogFilePath string

	RedisProtocol string
	RedisPort     string

	PostgresDsn string
	ExpiresSec  uint

	IsRelease bool
}

func InitSettings() (settings Settings) {
	settings = Settings{
		RootURL:    "/api",
		SessionURL: "/sessions",
		ProfileURL: "/profile",
		BoardsURL:  "/boards",

		ServerAddress: ":8000",

		Origins: []string{
			"http://localhost:8000",
			// Адрес деплоя
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		corsConfig: cors.DefaultConfig(),

		LogFilePath: "/var/log/backendLogs.log",

		RedisProtocol: "tcp",
		RedisPort:     ":6379",

		PostgresDsn: "host=localhost user=backend_ldbr password=backend_LDBR_password dbname=backend_ldbr_db port=5432 sslmode=disable",
		ExpiresSec:  60 * 24 * 3,

		IsRelease: false,
	}

	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	return
}
