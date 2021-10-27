package models

type Team struct {
	ID     uint    `json:"id"`
	Title  string  `json:"team_name" faker:"word"`
	Boards []Board `json:"boards" faker:"-"`
}
