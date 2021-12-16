package models

type User struct {
	UID            uint      `json:"uid" faker:"-" gorm:"primaryKey"`
	Login          string    `json:"login" faker:"username" gorm:"not null;unique;index"`
	Email          string    `json:"email" faker:"email" gorm:"not null;unique;index"`
	Password       string    `json:"password" faker:"password,len=10" gorm:"-"`
	OldPassword    string    `json:"old_password" faker:"password,len=10" gorm:"-"`
	HashedPassword []byte    `json:"-" faker:"-"`
	Description    string    `json:"description" faker:"sentence"`
	Avatar         string    `json:"avatar" faker:"uuid_digit"`
	Teams          []Team    `json:"teams" faker:"-" gorm:"many2many:users_teams;"`
	Boards         []Board   `json:"boards" faker:"-" gorm:"many2many:users_boards;"`
	AssignedCards  []Card    `json:"-" faker:"-" gorm:"many2many:users_cards;"`
	Comments       []Comment `json:"comments" gorm:"foreignKey:UID;constraint:OnDelete:CASCADE;"`
}
