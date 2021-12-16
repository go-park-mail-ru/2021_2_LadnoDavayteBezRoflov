package handlers

import (
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/sessionCookieController"
	"backendServer/pkg/tokens"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionMiddleware interface {
	CheckAuth() gin.HandlerFunc
	CSRF() gin.HandlerFunc
}

type SessionMiddlewareImpl struct {
	sessionUseCase usecases.SessionUseCase
}

func CreateSessionMiddleware(sessionUseCase usecases.SessionUseCase) SessionMiddleware {
	return &SessionMiddlewareImpl{sessionUseCase: sessionUseCase}
}

func (middleware *SessionMiddlewareImpl) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Request.Cookie("session_id")
		if err != nil {
			return
		}

		sid := session.Value
		uid, err := middleware.sessionUseCase.GetUID(sid)
		if err != nil {
			sessionCookieController.SetSessionCookieExpired(session)
			http.SetCookie(c.Writer, session)
			_ = c.Error(err)
			return
		}

		c.Set("uid", uid)
		c.Set("sid", sid)
	}
}

func tokenGetter(c *gin.Context) string {
	r := c.Request

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	} else {
		csrf, err := r.Cookie("csrf_token")
		if err == nil && len(csrf.Value) > 0 {
			return csrf.Value
		}
	}

	return ""
}

func (middleware *SessionMiddlewareImpl) CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, exists := c.Get("sid")
		if !exists {
			_ = c.Error(customErrors.ErrNotAuthorized)
			return
		}
		uid, exists := c.Get("uid")
		if !exists {
			_ = c.Error(customErrors.ErrNotAuthorized)
			return
		}

		jwtTokens := tokens.NewToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

		ignoreMethods := []string{"GET", "HEAD", "OPTIONS"}
		for _, method := range ignoreMethods {
			if c.Request.Method == method {
				token, err := jwtTokens.Create(sid.(string), uint32(uid.(uint)), time.Now().Add(15*time.Minute).Unix())
				if err != nil {
					_ = c.Error(err)
					return
				}
				c.Set("csrfToken", token)

				CSRFcookie := &http.Cookie{
					Name:     "csrf_token",
					Value:    token,
					Path:     "/",
					Expires:  time.Now().Add(15 * time.Minute),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				}

				http.SetCookie(c.Writer, CSRFcookie)
				c.Next()
				return
			}
		}

		_, err := jwtTokens.Check(sid.(string), uint32(uid.(uint)), tokenGetter(c))
		if err != nil {
			_ = c.Error(customErrors.ErrNoAccess)
			return
		}

		c.Next()
	}
}
