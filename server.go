package main

import (
	"backendServer/handlers"
	"backendServer/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"sync"
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
	router.Use(gin.Logger()) //TODO возможно заменить на другой логгер
	router.Use(gin.Recovery())
	router.Use(cors.New(server.settings.corsConfig))

	//TEMP DATA STORAGE
	data := &models.Data {
		Sessions: map[string]uint{},
		Users: map[string]*models.User {
			"Dima": {
				ID:       0,
				Login:    "Dima",
				Email:    "dima@mail.ru",
				Password: "ya_dima",
				Teams: []uint {
					0,
					1,
				},
			},
			"Tim": {
				ID:       1,
				Login:    "Tim",
				Email:    "tim@mail.ru",
				Password: "ya_tim",
			},
			"Tony": {
				ID:       2,
				Login:    "Tony",
				Email:    "tony@mail.ru",
				Password: "ya_tony",
				Teams: []uint{
					1,
				},
			},
		},
		Teams: map[uint]models.Team {
			0: {
				ID:     0,
				Title:  "Убийца Trello",
				Boards: []models.Board{
					{
						ID:    0,
						Title: "Убийца Trello",
						Description: "Неповторимый оригинал против жалкой пародии",
						Tasks: []string {
							"Начать делать",
							"Закончить",
						},
					},
					{
						ID:    1,
						Title: "Drello 2.0",
						Description: "Вторая часть знаменитого проекта",
						Tasks: []string{},
					},
					{
						ID:    2,
						Title: "Brrrello",
						Description: "Полностью оригинальная разработка",
						Tasks: []string {
							"Придумать оригинальное название",
						},
					},
				},
			},
			1: {
				ID:     1,
				Title:  "Сибирская корона 2021",
				Boards: []models.Board{
					{
						ID:    3,
						Title: "Технопарк",
						Description: "Здесь людей ...",
						Tasks: []string {
							"Не умереть",
							"Написать бэкэнд сервер",
						},
					},
					{
						ID:    4,
						Title: "Почилить",
						Description: "Планы по становлению овощем",
						Tasks: []string {
							"Даже не думай об этом",
							"Закрыть все таски из первой доски",
						},
					},
				},
			},
		},
		Mu:	&sync.RWMutex{},
	}

	rootGroup := router.Group(server.settings.RootURL)

	handlers.CreateSessionHandler(rootGroup, server.settings.SessionURL, data)
	handlers.CreateUserHandler(rootGroup, server.settings.ProfileURL, data)
	handlers.CreateBoardHandler(rootGroup, server.settings.BoardsURL, data)

	err := router.Run(server.settings.ServerAddress)
	if err != nil {
		panic(err) //TODO заменить на нормальную запись ошибки
	}
}
