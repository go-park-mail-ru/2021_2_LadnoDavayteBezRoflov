package models

type Board struct {
	BID         uint             `json:"bid" gorm:"primaryKey"`
	TID         uint             `json:"tid" gorm:"not null;index"`
	Title       string           `json:"board_name" faker:"word" gorm:"not null;index"`
	Description string           `json:"description" faker:"sentence"`
	Users       []User           `json:"-" faker:"-" gorm:"many2many:users_boards;"`
	Members     []PublicUserInfo `json:"members" faker:"-" gorm:"-"`
	CardLists   []CardList       `json:"card_lists" faker:"-" gorm:"foreignKey:BID;constraint:OnDelete:CASCADE;"`
	Cards       []Card           `json:"-" faker:"-" gorm:"foreignKey:BID;constraint:OnDelete:CASCADE;"`
}
