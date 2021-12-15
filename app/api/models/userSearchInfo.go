package models

//easyjson:json
type UsersSearchInfo []UserSearchInfo

//easyjson:json
type UserSearchInfo struct {
	UID    uint   `json:"uid"`
	Login  string `json:"userName"`
	Avatar string `json:"avatar"`
	Added  bool   `json:"added"`
}
