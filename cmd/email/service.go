package main

import (
	"backendServer/app/api/models"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/streadway/amqp"
)

type Service struct {
	settings Settings
}

func CreateService() *Service {
	settings := InitSettings()
	return &Service{settings: settings}
}

func (service *Service) Run() {
	// Logger and Closer
	var logger zapLogger.Logger
	logger.InitLogger(service.settings.LogFilePath)
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)

	// Mail
	mailDealer := gomail.NewDialer(
		service.settings.MailHost,
		service.settings.MailPort,
		service.settings.MailUsername,
		service.settings.MailPassword,
	)
	mailDealer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// RabbitMQ
	conn, err := amqp.Dial(service.settings.RabbitMQPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer everythingCloser.Close(conn.Close)

	channel, err := conn.Channel()
	if err != nil {
		logger.Error(err)
		return
	}
	defer everythingCloser.Close(channel.Close)

	queue, err := channel.QueueDeclare(
		service.settings.QueueName, // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		logger.Error(err)
		return
	}

	msgs, err := channel.Consume(
		queue.Name,                    // queue
		service.settings.ConsumerName, // consumer
		false,                         // auto-ack
		false,                         // exclusive
		false,                         // no-local
		false,                         // no-wait
		nil,                           // args
	)
	if err != nil {
		logger.Error(err)
		return
	}

	foreverChannel := make(chan bool)

	go func() {
		for d := range msgs {
			userInfo := new(models.PublicUserInfo)
			err = json.Unmarshal(d.Body, userInfo)
			if err != nil {
				logger.Error(err)
				break
			}
			fmt.Println(userInfo)
			emailLetter := gomail.NewMessage()
			emailLetter.SetHeader("From", service.settings.MailUsername)
			emailLetter.SetHeader("To", userInfo.Email)
			emailLetter.SetHeader("Subject", "Добро пожаловать в Brrrello!")
			emailLetter.SetBody("text/plain", "Рады вас видеть у себя на сайте!")
			if err = mailDealer.DialAndSend(emailLetter); err != nil {
				fmt.Println(err)
				logger.Error(err)
			}

			err = d.Ack(false)
			if err != nil {
				return
			}
		}
	}()

	<-foreverChannel
}
