package main

import (
	"fmt"
	"time"

	envParser "github.com/caarlos0/env"
	"github.com/gin-contrib/cors"

	"github.com/spf13/viper"
)

type Settings struct {
	RootURL           string `mapstructure:"root_url"`
	SessionURL        string `mapstructure:"session_url"`
	ProfileURL        string `mapstructure:"profile_url"`
	TeamsURL          string `mapstructure:"teams_url"`
	BoardsURL         string `mapstructure:"boards_url"`
	CardListsURL      string `mapstructure:"card_lists_url"`
	CardsURL          string `mapstructure:"cards_url"`
	TagsURL           string `mapstructure:"tags_url"`
	CommentsURL       string `mapstructure:"comments_url"`
	CheckListsURL     string `mapstructure:"check_lists_url"`
	CheckListItemsURL string `mapstructure:"check_list_items_url"`
	UserSearchURL     string `mapstructure:"user_search_url"`
	AttachmentsURL    string `mapstructure:"attachments_url"`

	ServerAddress         string `mapstructure:"server_address"`
	SessionServiceAddress string `mapstructure:"session_service_address"`
	RabbitMQAddress       string `mapstructure:"rabbitmq_address"`

	QueueName string `mapstructure:"queue_name"`

	Origins        []string
	AllowedMethods []string `mapstructure:"allowed_methods"`

	SessionCookieLifeTimeInDays time.Duration `mapstructure:"session_cookie_life_time_in_days"`

	corsConfig cors.Config

	LogFilePath       string
	AvatarsPath       string
	AttachmentsPath   string
	DefaultAvatarName string `mapstructure:"default_avatar_name"`

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

	viper.AddConfigPath("./cmd/api")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%+v\n", err)
	}

	settings = Settings{
		RootURL:           viper.GetString("url.root_url"),
		SessionURL:        viper.GetString("url.session_url"),
		ProfileURL:        viper.GetString("url.profile_url"),
		TeamsURL:          viper.GetString("url.teams_url"),
		BoardsURL:         viper.GetString("url.boards_url"),
		CardListsURL:      viper.GetString("url.card_lists_url"),
		CardsURL:          viper.GetString("url.cards_url"),
		TagsURL:           viper.GetString("url.tags_url"),
		CommentsURL:       viper.GetString("url.comments_url"),
		CheckListsURL:     viper.GetString("url.check_lists_url"),
		CheckListItemsURL: viper.GetString("url.check_list_items_url"),
		UserSearchURL:     viper.GetString("url.user_search_url"),
		AttachmentsURL:    viper.GetString("url.attachments_url"),

		ServerAddress:         viper.GetString("server_address"),
		SessionServiceAddress: viper.GetString("session_service_address"),
		RabbitMQAddress:       viper.GetString("rabbitmq_address"),

		QueueName: viper.GetString("queue_name"),

		Origins: []string{
			"http://localhost:8000",
			"http://prometheus:9090",
			fmt.Sprintf("https://%s", env.FRONTEND_ADDRESS),
		},

		AllowedMethods: viper.GetStringSlice("allowed_methods"),

		SessionCookieLifeTimeInDays: viper.GetDuration("session_cookie_life_time_in_days"),

		corsConfig: cors.DefaultConfig(),

		LogFilePath:       env.LOG_LOCATION,
		AvatarsPath:       env.FRONTEND_PATH,
		AttachmentsPath:   env.FRONTEND_PATH,
		DefaultAvatarName: viper.GetString("default_avatar_name"),

		PostgresDsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.DATABASE_HOST, env.POSTGRES_USER, env.POSTGRES_PASSWORD, env.POSTGRES_DB, env.DB_PORT),
	}
	settings.corsConfig.AllowOrigins = settings.Origins
	settings.corsConfig.AllowCredentials = true

	return
}
