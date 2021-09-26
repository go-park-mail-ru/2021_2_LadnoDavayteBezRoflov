package main

import (
	"backendServer/handlers"
	"backendServer/repositories"
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

	sessionRepo := repositories.CreateSessionRepository(data)
	userRepo := repositories.CreateUserRepository(data)
	boardRepo := repositories.CreateBoardRepository(data)

	rootGroup := router.Group(server.settings.RootURL)
	handlers.CreateSessionHandler(rootGroup, server.settings.SessionURL, sessionRepo)
	handlers.CreateUserHandler(rootGroup, server.settings.ProfileURL, userRepo, sessionRepo)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, boardRepo, sessionRepo)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		fmt.Println(err.Error())
	}
}
