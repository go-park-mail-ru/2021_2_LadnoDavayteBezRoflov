package handler

import (
	"backendServer/app/api/models"
	"backendServer/app/microservices/email/usecase"
	"backendServer/pkg/logger"
	"backendServer/pkg/metrics"
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
		metrics.EmailHits.WithLabelValues("500", "send first letter; can't consume").Inc()
		emailServer.logger.Error(err)
		return
	}

	foreverChannel := make(chan bool)

	go func() {
		for data := range messages {
			userInfo := new(models.PublicUserInfo)
			err := json.Unmarshal(data.Body, userInfo)
			if err != nil {
				metrics.EmailHits.WithLabelValues("500", "send first letter; can't unmarshal").Inc()
				emailServer.logger.Error(err)
				return
			}

			emailLetter := emailServer.EmailUseCase.SendFirstLetter(userInfo)
			if err = emailServer.mailDealer.DialAndSend(emailLetter); err != nil {
				metrics.EmailHits.WithLabelValues("500", "send first letter; can't send").Inc()
				emailServer.logger.Error(err)
				return
			}

			err = data.Ack(false)
			if err != nil {
				metrics.EmailHits.WithLabelValues("500", "send first letter; can't ack").Inc()
				emailServer.logger.Error(err)
				return
			}

			metrics.EmailHits.WithLabelValues("200", "send first letter").Inc()
		}
	}()

	<-foreverChannel
}

func (emailServer *EmailServer) SendNotifications() {
	for {
		go func() {
			emailLetters, err := emailServer.EmailUseCase.SendNotifications()
			if err != nil {
				metrics.EmailHits.WithLabelValues("500", "send notifications; can't get letters").Inc()
				emailServer.logger.Error(err)
			}

			for _, emailLetter := range *emailLetters {
				if err = emailServer.mailDealer.DialAndSend(&emailLetter); err != nil {
					metrics.EmailHits.WithLabelValues("500", "send notifications; can't send").Inc()
					emailServer.logger.Error(err)
					break
				}

				metrics.EmailHits.WithLabelValues("200", "send notifications").Inc()

				duration := time.Duration(30) * time.Minute
				time.Sleep(duration)
			}
		}()

		duration := time.Duration(24) * time.Hour
		time.Sleep(duration)
	}
}
