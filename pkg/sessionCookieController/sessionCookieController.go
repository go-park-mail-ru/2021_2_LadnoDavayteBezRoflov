package sessionCookieController

import (
	"net/http"
	"time"
)

var SessionCookieLifeTimeInHours time.Duration

func InitSessionCookieController(sessionCookieLifeTimeInDays time.Duration) {
	SessionCookieLifeTimeInHours = 24 * (sessionCookieLifeTimeInDays * time.Hour)
}

func CreateSessionCookie(sid string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		Expires:  time.Now().Add(SessionCookieLifeTimeInHours),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func SetSessionCookieExpired(cookie *http.Cookie) {
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(-1)
}
