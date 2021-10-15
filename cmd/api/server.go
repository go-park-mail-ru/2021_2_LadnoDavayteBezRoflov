package main

import (
	"backendServer/handlers"
	"backendServer/repositories/stores"
	useCases "backendServer/usecases/impl"
	"backendServer/utils"
	"fmt"
	"io"
	"os"

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
	logFile, logFileErr := os.Create(server.settings.LogFilePath)
	if logFileErr == nil {
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		fmt.Println(logFileErr.Error())
	}
	router.Use(server.settings.LogFormat)
	router.Use(gin.Recovery())
	router.Use(cors.New(server.settings.corsConfig))

	// TEMP DATA STORAGE
	data, err := utils.FillTestData(5, 3, 15)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sessionRepo := stores.CreateSessionRepository(data)
	userRepo := stores.CreateUserRepository(data)
	boardRepo := stores.CreateBoardRepository(data)

	sessionUseCase := useCases.CreateSessionUseCase(sessionRepo, userRepo)
	userUseCase := useCases.CreateUserUseCase(sessionRepo, userRepo)
	boardUseCase := useCases.CreateBoardUseCase(boardRepo)

	middleware := handlers.CreateMiddleware(sessionUseCase)

	rootGroup := router.Group(server.settings.RootURL)
	handlers.CreateSessionHandler(rootGroup, server.settings.SessionURL, sessionUseCase, middleware)
	handlers.CreateUserHandler(rootGroup, server.settings.ProfileURL, userUseCase)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, boardUseCase, middleware)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		fmt.Println(err.Error())
	}
}
