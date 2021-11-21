package models

type Team struct {
	TID    uint    `json:"tid" gorm:"primaryKey"`
	Title  string  `json:"team_name" faker:"word" gorm:"not null;unique;index"`
	Boards []Board `json:"boards" faker:"-" gorm:"foreignKey:TID;constraint:OnDelete:CASCADE;"`
	Users  []User  `json:"users" faker:"-" gorm:"many2many:users_teams;"`
}
