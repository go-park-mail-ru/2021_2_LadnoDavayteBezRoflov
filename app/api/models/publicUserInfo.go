package models

type PublicUserInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}
