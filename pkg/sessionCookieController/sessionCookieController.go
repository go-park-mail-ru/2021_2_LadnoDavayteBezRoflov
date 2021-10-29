package sessionCookieController

import (
	"net/http"
	"time"
)

var (
	SmallLifeTimeInHours         = 24 * time.Hour
	SessionCookieLifeTimeInDays  time.Duration
	SessionCookieLifeTimeInHours = 24 * SessionCookieLifeTimeInDays
	SessionCookieLifeTimeInSecs  = 60 * 60 * SessionCookieLifeTimeInHours
)

func InitSessionCookieController(sessionCookieLifeTimeInDays time.Duration) {
	SessionCookieLifeTimeInDays = sessionCookieLifeTimeInDays
}

func CreateSessionCookie(sid string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Expires:  time.Now().Add(time.Hour * SessionCookieLifeTimeInHours),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func SetSessionCookieExpired(cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(-1)
}

func UpdateSessionCookieExpires(cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(time.Hour * SessionCookieLifeTimeInHours)
}

func IsSessionCookieExpiresSoon(cookie *http.Cookie) bool {
	return time.Until(cookie.Expires).Hours() < SmallLifeTimeInHours.Hours()
}
