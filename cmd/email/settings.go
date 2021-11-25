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

	RabbitMQPath string
	QueueName    string
	ConsumerName string
}

type EnvironmentVariables struct {
	EMAIL_LOG_LOCATION string `env:"SESSION_LOG_LOCATION" envDefault:"/var/log/emailLogs.log"`
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
		MailPassword: viper.GetString("mail.password"),

		LogFilePath: env.EMAIL_LOG_LOCATION,

		RabbitMQPath: viper.GetString("rabbitmq_path"),
		QueueName:    viper.GetString("queue_name"),
		ConsumerName: viper.GetString("consumer_name"),
	}

	return
}
