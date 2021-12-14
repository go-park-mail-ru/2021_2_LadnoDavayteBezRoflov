package main

import (
	"backendServer/app/microservices/session/handler"
	handlerImpl "backendServer/app/microservices/session/handler/impl"
	"backendServer/app/microservices/session/repository/store"
	usecaseImpl "backendServer/app/microservices/session/usecase/impl"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"backendServer/pkg/metrics"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc/keepalive"

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

	// Prometheus metrics
	prometheus.MustRegister(metrics.SessionHits)

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

	// Prometheus server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Error(http.ListenAndServe(service.settings.ServiceMetricsPort, nil))
	}()

	// GRPC server
	grpcSrv := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
	)
	handler.RegisterSessionCheckerServer(grpcSrv, handlerImpl.CreateSessionCheckerServer(sessionUseCase))
	if err = grpcSrv.Serve(listener); err != nil {
		logger.Error(err)
	}
}
