package models

type PublicUserInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName" faker:"username,unique"`
	Avatar string `json:"avatar"`
}
