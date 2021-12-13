package handler

import (
	"backendServer/app/api/models"
	"backendServer/app/microservices/email/usecase"
	"backendServer/pkg/logger"
	"encoding/json"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/streadway/amqp"
)

type EmailServer struct {
	EmailUseCase usecase.EmailUseCase
	logger       *logger.Logger
	mailDealer   *gomail.Dialer
	channel      *amqp.Channel
	queueName    string
	consumerName string
}

func CreateEmailServer(
	emailUseCase usecase.EmailUseCase,
	logger *logger.Logger,
	mailDealer *gomail.Dialer,
	channel *amqp.Channel,
	queueName string,
	consumerName string,
) *EmailServer {
	return &EmailServer{
		EmailUseCase: emailUseCase,
		logger:       logger,
		mailDealer:   mailDealer,
		channel:      channel,
		queueName:    queueName,
		consumerName: consumerName,
	}
}

func (emailServer *EmailServer) Run() {
	go func() {
		emailServer.SendFirstLetter(emailServer.channel)
	}()
	emailServer.SendNotifications()
}

func (emailServer *EmailServer) SendFirstLetter(channel *amqp.Channel) {
	messages, err := channel.Consume(
		emailServer.queueName,    // queue
		emailServer.consumerName, // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		emailServer.logger.Error(err)
		return
	}

	foreverChannel := make(chan bool)

	go func() {
		for data := range messages {
			userInfo := new(models.PublicUserInfo)
			err := json.Unmarshal(data.Body, userInfo)
			if err != nil {
				emailServer.logger.Error(err)
				return
			}

			emailLetter := emailServer.EmailUseCase.SendFirstLetter(userInfo)
			if err = emailServer.mailDealer.DialAndSend(emailLetter); err != nil {
				emailServer.logger.Error(err)
				return
			}

			err = data.Ack(false)
			if err != nil {
				emailServer.logger.Error(err)
				return
			}
		}
	}()

	<-foreverChannel
}

func (emailServer *EmailServer) SendNotifications() {
	for true {
		go func() {
			emailLetters, err := emailServer.EmailUseCase.SendNotifications()
			if err != nil {
				emailServer.logger.Error(err)
			}

			for _, emailLetter := range *emailLetters {
				if err = emailServer.mailDealer.DialAndSend(&emailLetter); err != nil {
					emailServer.logger.Error(err)
					break
				}

				duration := time.Duration(30) * time.Minute
				time.Sleep(duration)
			}
		}()

		duration := time.Duration(12) * time.Hour
		time.Sleep(duration)
	}
}
