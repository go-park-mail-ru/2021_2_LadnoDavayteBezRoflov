package main

import (
	"backendServer/app/api/models"
	"backendServer/app/microservices/email/handler"
	"backendServer/app/microservices/email/repository/store"
	"backendServer/app/microservices/email/usecase/impl"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"backendServer/pkg/metrics"
	"crypto/tls"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gopkg.in/gomail.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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
	logger := new(zapLogger.Logger)
	logger.InitLogger(service.settings.LogFilePath)
	everythingCloser := closer.CreateCloser(logger)
	defer everythingCloser.Close(logger.Sync)

	// Postgres
	postgresClient, err := gorm.Open(postgres.Open(service.settings.PostgresDsn), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		return
	}
	err = postgresClient.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Board{},
		&models.CardList{},
		&models.Card{},
		&models.Comment{},
		&models.CheckList{},
		&models.CheckListItem{},
	)
	if err != nil {
		logger.Error(err)
		return
	}

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

	// Prometheus metrics
	prometheus.MustRegister(metrics.EmailHits)

	// Repository
	emailRepo := store.CreateEmailRepository(postgresClient)

	// UseCase
	emailUseCase := impl.CreateEmailUseCase(emailRepo, service.settings.MailUsername)

	// Handler
	emailHandler := handler.CreateEmailServer(emailUseCase, logger, mailDealer, channel, queue.Name, service.settings.ConsumerName)

	// Prometheus server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Error(http.ListenAndServe(service.settings.ServiceMetricsPort, nil))
	}()

	emailHandler.Run()
}
