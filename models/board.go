package models

type Board struct {
	ID          uint     `json:"id"`
	Title       string   `json:"board_name" faker:"word"`
	Description string   `json:"description" faker:"sentence"`
	Tasks       []string `json:"tasks" faker:"sentence,slice_len=3"` // TODO обернуть в нормальные задачи (позже)
}
