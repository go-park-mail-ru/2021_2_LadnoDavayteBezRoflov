package main

import (
	"fmt"

	envParser "github.com/caarlos0/env"
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

	settings = Settings{
		MailHost:     "smtp.mail.ru",
		MailPort:     587,
		MailUsername: "brrrello-notify@mail.ru",
		MailPassword: "B1GmLskQvzhubYNKxKq0",

		LogFilePath: env.EMAIL_LOG_LOCATION,

		RabbitMQPath: "amqp://guest:guest@rabbitmq:5672/",
		QueueName:    "queue",
		ConsumerName: "consumer",
	}

	return
}
