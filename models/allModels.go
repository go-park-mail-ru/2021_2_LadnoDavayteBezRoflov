package models

import "sync"

type Board struct {
	ID		uint `json:"id"`
	Title	string `json:"title"`
	Tasks	[]string `json:"tasks"` //TODO обернуть в нормальные задачи (позже)
}

type User struct {
	ID			uint `json:"id"`
	Login		string `json:"login"`
	Email		string `json:"email"`
	Password	string `json:"password"`
}

//TEMP DATA STORAGE MODEL
type Data struct {
	Sessions map[string]uint
	Users    map[string]*User
	Boards   map[uint][]Board
	Mu    *sync.RWMutex
}
