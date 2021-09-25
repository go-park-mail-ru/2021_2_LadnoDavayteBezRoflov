package models

import "sync"

type Board struct {
	ID		uint `json:"id"`
	Title	string `json:"board_name"`
	Tasks	[]string `json:"tasks"` //TODO обернуть в нормальные задачи (позже)
}

type Team struct {
	ID		uint `json:"id"`
	Title	string `json:"team_name"`
	Boards	[]Board `json:"boards"`
}

type User struct {
	ID			uint `json:"id"`
	Login		string `json:"login"`
	Email		string `json:"email"`
	Password	string `json:"password"`
	Teams		[]uint `json:"teams"`
}

//TEMP DATA STORAGE MODEL
type Data struct {
	Sessions	map[string]uint
	Users		map[string]*User
	//Boards		map[uint][]Board
	Teams		map[uint]Team
	Mu			*sync.RWMutex
}
