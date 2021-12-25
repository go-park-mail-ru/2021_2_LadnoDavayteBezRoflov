package main

import (
	"backendServer/app/api/handlers"
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/stores"
	"backendServer/app/api/usecases/impl"
	"backendServer/app/microservices/session/handler"
	"backendServer/pkg/closer"
	zapLogger "backendServer/pkg/logger"
	"backendServer/pkg/metrics"
	"backendServer/pkg/sessionCookieController"
	"backendServer/pkg/webSockets"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/gin-contrib/expvar"

	"github.com/streadway/amqp"

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
		&models.Attachment{},
		&models.Tag{},
	)
	if err != nil {
		logger.Error(err)
		return
	}

	// RabbitMQ
	conn, err := amqp.Dial(server.settings.RabbitMQAddress)
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

	_, err = channel.QueueDeclare(
		server.settings.QueueName, // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
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
	userRepo := stores.CreateUserRepository(postgresClient, server.settings.AvatarsPath, server.settings.DefaultAvatarName, channel, server.settings.QueueName)
	teamRepo := stores.CreateTeamRepository(postgresClient)
	boardRepo := stores.CreateBoardRepository(postgresClient)
	cardListRepo := stores.CreateCardListRepository(postgresClient)
	cardRepo := stores.CreateCardRepository(postgresClient)
	tagRepo := stores.CreateTagRepository(postgresClient)
	commentRepo := stores.CreateCommentRepository(postgresClient)
	checkListRepo := stores.CreateCheckListRepository(postgresClient)
	checkListItemRepo := stores.CreateCheckListItemRepository(postgresClient)
	attachmentRepo := stores.CreateAttachmentRepository(postgresClient, server.settings.AttachmentsPath)

	// UseCases
	sessionUseCase := impl.CreateSessionUseCase(sessionRepo, userRepo)
	userUseCase := impl.CreateUserUseCase(sessionRepo, userRepo, teamRepo)
	teamUseCase := impl.CreateTeamUseCase(teamRepo, userRepo, boardRepo)
	boardUseCase := impl.CreateBoardUseCase(boardRepo, userRepo, teamRepo, cardListRepo, cardRepo, checkListRepo)
	cardListUseCase := impl.CreateCardListUseCase(cardListRepo, userRepo)
	cardUseCase := impl.CreateCardUseCase(cardRepo, userRepo, tagRepo)
	commentUseCase := impl.CreateCommentUseCase(commentRepo, userRepo)
	tagUseCase := impl.CreateTagUseCase(tagRepo, userRepo)
	checkListUseCase := impl.CreateCheckListUseCase(checkListRepo, userRepo)
	checkListItemUseCase := impl.CreateCheckListItemUseCase(checkListItemRepo, userRepo)
	userSearchUseCase := impl.CreateUserSearchUseCase(userRepo, cardRepo, teamRepo, boardRepo)
	attachmentUseCase := impl.CreateAttachmentUseCase(attachmentRepo, userRepo)

	// Middlewares
	commonMiddleware := handlers.CreateCommonMiddleware(logger)
	sessionMiddleware := handlers.CreateSessionMiddleware(sessionUseCase)

	router.Use(commonMiddleware.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(server.settings.corsConfig))

	// get global Monitor object
	monitor := ginmetrics.GetMonitor()
	_ = ginmetrics.GetMonitor().AddMetric(metrics.APIErrors)
	monitor.SetMetricPath("/metrics")
	monitor.SetSlowTime(10)
	monitor.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	monitor.Use(router)

	webSockets.SetupWebSocketHandler()

	// Handlers
	router.NoRoute(handlers.NoRouteHandler)
	router.GET("/debug/vars", expvar.Handler())
	router.GET("/ws", sessionMiddleware.CheckAuth(), sessionMiddleware.CSRF(), webSockets.WebSocketsHandler)
	rootGroup := router.Group(server.settings.RootURL)
	handlers.CreateSessionHandler(rootGroup, server.settings.SessionURL, sessionUseCase, sessionMiddleware)
	handlers.CreateUserHandler(rootGroup, server.settings.ProfileURL, userUseCase, sessionMiddleware)
	handlers.CreateTeamHandler(rootGroup, server.settings.TeamsURL, teamUseCase, sessionMiddleware)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, boardUseCase, sessionMiddleware)
	handlers.CreateCardListHandler(rootGroup, server.settings.CardListsURL, cardListUseCase, sessionMiddleware)
	handlers.CreateCardHandler(rootGroup, server.settings.CardsURL, cardUseCase, sessionMiddleware)
	handlers.CreateCommentHandler(rootGroup, server.settings.CommentsURL, commentUseCase, sessionMiddleware)
	handlers.CreateTagHandler(rootGroup, server.settings.TagsURL, tagUseCase, sessionMiddleware)
	handlers.CreateCheckListHandler(rootGroup, server.settings.CheckListsURL, checkListUseCase, sessionMiddleware)
	handlers.CreateCheckListItemHandler(rootGroup, server.settings.CheckListItemsURL, checkListItemUseCase, sessionMiddleware)
	handlers.CreateUserSearchHandler(rootGroup, server.settings.UserSearchURL, userSearchUseCase, sessionMiddleware)
	handlers.CreateAttachmentHandler(rootGroup, server.settings.AttachmentsURL, attachmentUseCase, sessionMiddleware)

	err = router.Run(server.settings.ServerAddress)
	if err != nil {
		logger.Error(err)
	}
}
