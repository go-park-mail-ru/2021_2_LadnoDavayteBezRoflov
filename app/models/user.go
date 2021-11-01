package models

type User struct {
	UID            uint   `json:"uid" gorm:"primaryKey"`
	Login          string `json:"login" faker:"username,unique" gorm:"not null;unique;index"`
	Email          string `json:"email" faker:"email,unique" gorm:"not null;unique;index"`
	Password       string `json:"password" faker:"password,len=10" gorm:"-"`
	HashedPassword []byte `json:"-" faker:"-"`
	Description    string `json:"description" faker:"sentence"`
	Avatar         string `json:"avatar" faker:"-"`
	Teams          []Team `json:"teams" faker:"-" gorm:"many2many:users_teams;"`
}
