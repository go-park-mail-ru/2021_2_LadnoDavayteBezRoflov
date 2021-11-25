package models

type PublicUserInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName"`
	Avatar string `json:"avatar"`
}
