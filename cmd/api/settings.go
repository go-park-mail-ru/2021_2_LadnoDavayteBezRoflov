package main

import (
	"fmt"
	"time"

	envParser "github.com/caarlos0/env"
	"github.com/gin-contrib/cors"
)

type Settings struct {
	RootURL           string
	SessionURL        string
	ProfileURL        string
	TeamsURL          string
	BoardsURL         string
	CardListsURL      string
	CardsURL          string
	CommentsURL       string
	CheckListsURL     string
	CheckListItemsURL string
	UserSearchURL     string

	ServerAddress         string
	SessionServiceAddress string
	RabbitMQAddress       string

	QueueName string

	Origins        []string
	AllowedMethods []string

	SessionCookieLifeTimeInDays time.Duration

	corsConfig cors.Config

	LogFilePath       string
	AvatarsPath       string
	DefaultAvatarName string

	PostgresDsn string
}

type EnvironmentVariables struct {
	DB_PORT           string `env:"DB_PORT,required"`
	POSTGRES_USER     string `env:"POSTGRES_USER,required"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD,required"`
	DATABASE_HOST     string `env:"DATABASE_HOST,required"`
	POSTGRES_DB       string `env:"POSTGRES_DB,required"`
	FRONTEND_ADDRESS  string `env:"FRONTEND_ADDRESS,required"`
	FRONTEND_PATH     string `env:"PUBLIC_DIR,required"`
	LOG_LOCATION      string `env:"LOG_LOCATION" envDefault:"/var/log/backendLogs.log"`
}

func InitSettings() (settings Settings) {
	env := EnvironmentVariables{}
	if err := envParser.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	settings = Settings{
		RootURL:           "/api",
		SessionURL:        "/sessions",
		ProfileURL:        "/profile",
		TeamsURL:          "/teams",
		BoardsURL:         "/boards",
		CardListsURL:      "/cardLists",
		CardsURL:          "/cards",
		CommentsURL:       "/comments",
		CheckListsURL:     "/checkLists",
		CheckListItemsURL: "/checkListItems",
		UserSearchURL:     "/usersearch",

		ServerAddress:         ":8000",
		SessionServiceAddress: "session:8081",
		RabbitMQAddress:       "amqp://guest:guest@rabbitmq:5672/",

		QueueName: "queue",

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

		LogFilePath:       env.LOG_LOCATION,
		AvatarsPath:       env.FRONTEND_PATH,
		DefaultAvatarName: "default_user_picture.webp",

		PostgresDsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.DATABASE_HOST, env.POSTGRES_USER, env.POSTGRES_PASSWORD, env.POSTGRES_DB, env.DB_PORT),
	}

	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	return
}
