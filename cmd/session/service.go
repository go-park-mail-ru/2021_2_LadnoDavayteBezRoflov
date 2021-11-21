package main

import (
	handler "backendServer/app/microservices/session/handler"
	handlerImpl "backendServer/app/microservices/session/handler/impl"
	"backendServer/app/microservices/session/repository/store"
	usecaseImpl "backendServer/app/microservices/session/usecase/impl"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/gomodule/redigo/redis"
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

	// Redis
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(service.settings.RedisProtocol, service.settings.RedisPort)
			if err != nil {
				panic(err)
			}
			return c, err
		},
	}
	defer everythingCloser.Close(redisPool.Close)

	// Repository
	sessionRepo := store.CreateSessionRepository(redisPool, uint64(24*(3*time.Hour)), everythingCloser)

	// UseCase
	sessionUseCase := usecaseImpl.CreateSessionUseCase(sessionRepo)

	// Handler
	listener, err := net.Listen(service.settings.ServiceProtocol, service.settings.ServicePort)
	if err != nil {
		logger.Error(err)
		return
	}
	defer everythingCloser.Close(listener.Close)

	grpcSrv := grpc.NewServer()
	handler.RegisterSessionCheckerServer(grpcSrv, handlerImpl.CreateSessionCheckerServer(sessionUseCase))
	if err = grpcSrv.Serve(listener); err != nil {
		logger.Error(err)
	}
}
