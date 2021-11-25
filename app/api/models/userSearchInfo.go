package models

type UserSearchInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName"`
	Avatar string `json:"avatar"`
	Added  bool   `json:"added"`
}
