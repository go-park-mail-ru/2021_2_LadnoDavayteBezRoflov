package main

import (
	"backendServer/app/handlers"
	"backendServer/app/models"
	"backendServer/app/repositories/stores"
	"backendServer/app/usecases/impl"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"backendServer/pkg/sessionCookieController"

	"github.com/gomodule/redigo/redis"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	settings Settings
}

func CreateServer() *Server {
	settings := InitSettings()
	sessionCookieController.InitSessionCookieController(settings.SessionCookieLifeTimeInDays)
	return &Server{settings: settings}
}

func (server *Server) Run() {
	router := gin.New()

	// Logger and Closer
	var logger zapLogger.Logger
	logger.InitLogger(server.settings.LogFilePath)
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)

	// Redis
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(server.settings.RedisProtocol, server.settings.RedisPort)
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			return c, err
		},
	}
	defer everythingCloser.Close(redisPool.Close)

	// Postgres
	postgresClient, err := gorm.Open(postgres.Open(server.settings.PostgresDsn), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		return
	}
	err = postgresClient.AutoMigrate(&models.User{}, &models.Team{}, &models.Board{}, &models.CardList{}, &models.Card{})
	if err != nil {
		logger.Error(err)
		return
	}

	// Repositories
	sessionRepo := stores.CreateSessionRepository(redisPool, uint64(sessionCookieController.SessionCookieLifeTimeInHours), everythingCloser)
	userRepo := stores.CreateUserRepository(postgresClient, server.settings.AvatarsPath, server.settings.DefaultAvatarName)
	teamRepo := stores.CreateTeamRepository(postgresClient)
	boardRepo := stores.CreateBoardRepository(postgresClient)
	cardListRepo := stores.CreateCardListRepository(postgresClient)
	cardRepo := stores.CreateCardRepository(postgresClient)
	commentRepo := stores.CreateCommentRepository(postgresClient)
	checkListRepo := stores.CreateCheckListRepository(postgresClient)
	checkListItemRepo := stores.CreateCheckListItemRepository(postgresClient)

	// UseCases
	sessionUseCase := impl.CreateSessionUseCase(sessionRepo, userRepo)
	userUseCase := impl.CreateUserUseCase(sessionRepo, userRepo, teamRepo)
	boardUseCase := impl.CreateBoardUseCase(boardRepo, userRepo, teamRepo, cardListRepo, cardRepo, checkListRepo)
	cardListUseCase := impl.CreateCardListUseCase(cardListRepo, userRepo)
	cardUseCase := impl.CreateCardUseCase(cardRepo, userRepo)
	commentUseCase := impl.CreateCommentUseCase(commentRepo, userRepo)
	checkListUseCase := impl.CreateCheckListUseCase(checkListRepo, userRepo)
	checkListItemUseCase := impl.CreateCheckListItemUseCase(checkListItemRepo, userRepo)

	// Middlewares
	commonMiddleware := handlers.CreateCommonMiddleware(logger)
	sessionMiddleware := handlers.CreateSessionMiddleware(sessionUseCase)

	router.Use(commonMiddleware.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(server.settings.corsConfig))

	// Handlers
	router.NoRoute(handlers.NoRouteHandler)
	rootGroup := router.Group(server.settings.RootURL)
	handlers.CreateSessionHandler(rootGroup, server.settings.SessionURL, sessionUseCase, sessionMiddleware)
	handlers.CreateUserHandler(rootGroup, server.settings.ProfileURL, userUseCase, sessionMiddleware)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, boardUseCase, sessionMiddleware)
	handlers.CreateCardListHandler(rootGroup, server.settings.CardListsURL, cardListUseCase, sessionMiddleware)
	handlers.CreateCardHandler(rootGroup, server.settings.CardsURL, cardUseCase, sessionMiddleware)
	handlers.CreateCommentHandler(rootGroup, server.settings.CommentsURL, commentUseCase, sessionMiddleware)
	handlers.CreateCheckListHandler(rootGroup, server.settings.CheckListsURL, checkListUseCase, sessionMiddleware)
	handlers.CreateCheckListItemHandler(rootGroup, server.settings.CheckListItemsURL, checkListItemUseCase, sessionMiddleware)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		logger.Error(err)
	}
}
