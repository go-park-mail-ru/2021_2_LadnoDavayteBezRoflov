package main

import (
	"fmt"

	envParser "github.com/caarlos0/env"
	"github.com/spf13/viper"
)

type Settings struct {
	MailHost     string
	MailPort     int
	MailUsername string
	MailPassword string

	LogFilePath string

	ServiceMetricsPort string

	RabbitMQPath string
	QueueName    string
	ConsumerName string

	PostgresDsn string
}

type EnvironmentVariables struct {
	EMAIL_LOG_LOCATION string `env:"EMAIL_LOG_LOCATION" envDefault:"/var/log/emailLogs.log"`
	EMAIL_PASSWORD     string `env:"EMAIL_PASSWORD,required"`
	DB_PORT            string `env:"DB_PORT,required"`
	POSTGRES_USER      string `env:"POSTGRES_USER,required"`
	POSTGRES_PASSWORD  string `env:"POSTGRES_PASSWORD,required"`
	DATABASE_HOST      string `env:"DATABASE_HOST,required"`
	POSTGRES_DB        string `env:"POSTGRES_DB,required"`
}

func InitSettings() (settings Settings) {
	env := EnvironmentVariables{}
	if err := envParser.Parse(&env); err != nil {
		fmt.Printf("%+v\n", err)
	}

	viper.AddConfigPath("./cmd/email")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%+v\n", err)
	}

	settings = Settings{
		MailHost:     viper.GetString("mail.host"),
		MailPort:     viper.GetInt("mail.port"),
		MailUsername: viper.GetString("mail.username"),
		MailPassword: env.EMAIL_PASSWORD,

		LogFilePath: env.EMAIL_LOG_LOCATION,

		ServiceMetricsPort: viper.GetString("service_metrics_port"),

		RabbitMQPath: viper.GetString("rabbitmq_path"),
		QueueName:    viper.GetString("queue_name"),
		ConsumerName: viper.GetString("consumer_name"),

		PostgresDsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.DATABASE_HOST, env.POSTGRES_USER, env.POSTGRES_PASSWORD, env.POSTGRES_DB, env.DB_PORT),
	}

	return
}
