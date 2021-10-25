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

		LogFilePath: "../../backendLogs.log",

		IsRelease: false,
	}

	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	return
}
