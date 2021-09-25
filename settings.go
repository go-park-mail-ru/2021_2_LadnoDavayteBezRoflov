package main

import "github.com/gin-contrib/cors"

type Settings struct {
	//TODO
	RootURL			string
	SessionURL		string
	ProfileURL		string
	BoardsURL		string

	ServerAddress	string

	Origins			[]string

	corsConfig		cors.Config
}

func InitSettings() (settings Settings) {
	//TODO
	settings = Settings{
		RootURL: "/api",
		SessionURL: "/sessions",
		ProfileURL: "/profile",
		BoardsURL: "/boards",

		ServerAddress: "0.0.0.0:8080",

		Origins: []string{
			"http://localhost:8080",
			//"", //Адрес деплоя
		},

		corsConfig: cors.DefaultConfig(),
	}

	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	return
}