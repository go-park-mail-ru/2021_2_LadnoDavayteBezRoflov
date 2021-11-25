package main

import (
	"backendServer/app/api/handlers"
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/stores"
	"backendServer/app/api/usecases/impl"
	"backendServer/app/microservices/session/handler"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"backendServer/pkg/sessionCookieController"

	"google.golang.org/grpc"

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

	// Postgres
	postgresClient, err := gorm.Open(postgres.Open(server.settings.PostgresDsn), &gorm.Config{})
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

	// Session microservice
	grpcConn, err := grpc.Dial(
		server.settings.SessionServiceAddress,
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Error(err)
		return
	}
	defer everythingCloser.Close(grpcConn.Close)

	sessionManager := handler.NewSessionCheckerClient(grpcConn)

	// Repositories
	sessionRepo := stores.CreateSessionRepository(sessionManager)
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
	teamUseCase := impl.CreateTeamUseCase(teamRepo, userRepo)
	boardUseCase := impl.CreateBoardUseCase(boardRepo, userRepo, teamRepo, cardListRepo, cardRepo, checkListRepo)
	cardListUseCase := impl.CreateCardListUseCase(cardListRepo, userRepo)
	cardUseCase := impl.CreateCardUseCase(cardRepo, userRepo)
	commentUseCase := impl.CreateCommentUseCase(commentRepo, userRepo)
	checkListUseCase := impl.CreateCheckListUseCase(checkListRepo, userRepo)
	checkListItemUseCase := impl.CreateCheckListItemUseCase(checkListItemRepo, userRepo)
	userSearchUseCase := impl.CreateUserSearchUseCase(userRepo, cardRepo, teamRepo, boardRepo)

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
	handlers.CreateTeamHandler(rootGroup, server.settings.TeamsURL, teamUseCase, sessionMiddleware)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, boardUseCase, sessionMiddleware)
	handlers.CreateCardListHandler(rootGroup, server.settings.CardListsURL, cardListUseCase, sessionMiddleware)
	handlers.CreateCardHandler(rootGroup, server.settings.CardsURL, cardUseCase, sessionMiddleware)
	handlers.CreateCommentHandler(rootGroup, server.settings.CommentsURL, commentUseCase, sessionMiddleware)
	handlers.CreateCheckListHandler(rootGroup, server.settings.CheckListsURL, checkListUseCase, sessionMiddleware)
	handlers.CreateCheckListItemHandler(rootGroup, server.settings.CheckListItemsURL, checkListItemUseCase, sessionMiddleware)
	handlers.CreateUserSearchHandler(rootGroup, server.settings.UserSearchURL, userSearchUseCase, sessionMiddleware)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		logger.Error(err)
	}
}
