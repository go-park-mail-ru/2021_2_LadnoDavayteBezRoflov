package models

type PublicUserInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName" faker:"username"`
	Email  string `json:"email" faker:"email"`
	Avatar string `json:"avatar"`
}
