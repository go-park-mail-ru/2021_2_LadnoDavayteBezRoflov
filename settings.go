package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Settings struct {
	// TODO
	RootURL    string
	SessionURL string
	ProfileURL string
	BoardsURL  string

	ServerAddress string

	Origins        []string
	AllowedMethods []string

	corsConfig cors.Config

	LogFilePath string
	LogFormat   gin.HandlerFunc
}

func InitSettings() (settings Settings) {
	settings = Settings{
		RootURL:    "/api",
		SessionURL: "/sessions",
		ProfileURL: "/profile",
		BoardsURL:  "/boards",

		ServerAddress: "0.0.0.0:8080",

		Origins: []string{
			"http://localhost:8080",
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

		LogFilePath: "backendLogs.log",
	}

	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	settings.LogFormat = gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})

	return
}
