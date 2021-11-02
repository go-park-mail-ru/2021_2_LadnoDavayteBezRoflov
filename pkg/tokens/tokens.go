package tokens

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Secret []byte
}

type JwtCsrfClaims struct {
	SessionID string
	UserID    uint32
	jwt.StandardClaims
}

func NewToken(secret string) *Token {
	return &Token{Secret: []byte(secret)}
}

func (token *Token) Create(sid string, uid uint32, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		SessionID: sid,
		UserID:    uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return jwtToken.SignedString(token.Secret)
}

func (token *Token) parseSecretGetter(tk *jwt.Token) (interface{}, error) {
	method, ok := tk.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return token.Secret, nil
}

func (token *Token) Check(sid string, uid uint32, inputToken string) (bool, error) {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, token.parseSecretGetter)
	if err != nil {
		return false, fmt.Errorf("cant parse jwt token: %v", err)
	}
	if payload.Valid() != nil {
		return false, fmt.Errorf("invalid jwt token: %v", err)
	}
	return payload.SessionID == sid && payload.UserID == uid, nil
}
