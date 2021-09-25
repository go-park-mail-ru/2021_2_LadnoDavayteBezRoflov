package handlers

import (
	"backendServer/models"
	"backendServer/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserHandler struct {
	UserURL	string
	Data		*models.Data
}

func CreateUserHandler(router *gin.RouterGroup, userURL string, data *models.Data) {
	handler := &UserHandler{
		UserURL:	userURL,
		Data:		data,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("", handler.Create)
	}
}

func (sessionHandler *UserHandler) Create(c *gin.Context) {
	var json models.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Bad request")})
		return
	}

	//TODO валидация данных

	sessionHandler.Data.Mu.RLock()
	_, userAlreadyCreated := sessionHandler.Data.Users[json.Login]
	users := sessionHandler.Data.Users
	sessionHandler.Data.Mu.RUnlock()

	if userAlreadyCreated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Bad input data")}) //TODO заменить на уникальный тип ошибки
		return
	}

	for _, user := range users {
		if user.Email == json.Email {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Bad input data")}) //TODO заменить на уникальный тип ошибки
			return
		}
	}

	userID := uint(len(users))
	sessionHandler.Data.Mu.Lock()
	sessionHandler.Data.Users[json.Login] = &models.User{
		ID: userID,
		Login:    json.Login,
		Email:    json.Email,
		Password: json.Password,
	}
	sessionHandler.Data.Mu.Unlock()

	SID := utils.RandStringRunes(32)

	sessionHandler.Data.Mu.Lock()
	sessionHandler.Data.Sessions[SID] = userID
	sessionHandler.Data.Mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(24 * time.Hour),
		Secure: true,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}
