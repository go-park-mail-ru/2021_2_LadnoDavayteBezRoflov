package models

type User struct {
	ID       uint   `json:"id"`
	Login    string `json:"login" faker:"username,unique"`
	Email    string `json:"email" faker:"email,unique"`
	Password string `json:"password" faker:"password,len=10"`
	Teams    []uint `json:"teams" faker:"-"`
}