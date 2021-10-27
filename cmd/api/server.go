package main

import (
	"backendServer/app/handlers"
	"backendServer/app/repositories/stores"
	"backendServer/app/usecases/impl"
	zapLogger "backendServer/pkg/logger"
	"backendServer/pkg/utils"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	settings Settings
}

func CreateServer() *Server {
	settings := InitSettings()
	return &Server{settings: settings}
}

func (server *Server) Run() {
	if server.settings.IsRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	var logger zapLogger.Logger
	logger.InitLogger(server.settings.LogFilePath)
	defer func(err error) {
		if err != nil {
			fmt.Println(err)
		}
	}(logger.Sync())

	// TEMP DATA STORAGE
	data, err := utils.FillTestData(5, 3, 15)
	if err != nil {
		logger.Error(err)
		return
	}

	sessionRepo := stores.CreateSessionRepository(data)
	userRepo := stores.CreateUserRepository(data)
	boardRepo := stores.CreateBoardRepository(data)

	sessionUseCase := impl.CreateSessionUseCase(sessionRepo, userRepo)
	userUseCase := impl.CreateUserUseCase(sessionRepo, userRepo)
	boardUseCase := impl.CreateBoardUseCase(boardRepo)

	commonMiddleware := handlers.CreateCommonMiddleware(logger)
	sessionMiddleware := handlers.CreateSessionMiddleware(sessionUseCase)

	router.Use(commonMiddleware.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(server.settings.corsConfig))

	rootGroup := router.Group(server.settings.RootURL)
	handlers.CreateSessionHandler(rootGroup, server.settings.SessionURL, sessionUseCase, sessionMiddleware)
	handlers.CreateUserHandler(rootGroup, server.settings.ProfileURL, userUseCase)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, boardUseCase, sessionMiddleware)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		logger.Error(err)
	}
}
