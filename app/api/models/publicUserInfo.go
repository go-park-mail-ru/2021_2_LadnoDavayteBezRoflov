package models

type PublicUserInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName" faker:"username,unique"`
	Email  string `json:"email" faker:"email,unique"`
	Avatar string `json:"avatar"`
}
