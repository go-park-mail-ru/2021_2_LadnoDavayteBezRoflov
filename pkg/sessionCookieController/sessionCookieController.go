package sessionCookieController

import (
	"net/http"
	"time"
)

var (
	SmallLifeTimeInHours         = 24 * time.Hour
	SessionCookieLifeTimeInHours time.Duration
)

func InitSessionCookieController(sessionCookieLifeTimeInDays time.Duration) {
	SessionCookieLifeTimeInHours = 24 * (sessionCookieLifeTimeInDays * time.Hour)
}

func CreateSessionCookie(sid string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Expires:  time.Now().Add(SessionCookieLifeTimeInHours),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func SetSessionCookieExpired(cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(-1)
}

func UpdateSessionCookieExpires(cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(SessionCookieLifeTimeInHours)
}

func IsSessionCookieExpiresSoon(cookie *http.Cookie) bool {
	return time.Until(cookie.Expires) < SmallLifeTimeInHours
}
