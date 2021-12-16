package models

type Board struct {
	BID             uint             `json:"bid" gorm:"primaryKey"`
	TID             uint             `json:"tid" gorm:"not null;index"`
	Title           string           `json:"board_name" faker:"word" gorm:"not null;index"`
	Description     string           `json:"description" faker:"sentence"`
	AccessPath      string           `json:"access_path"`
	Users           []User           `json:"-" gorm:"many2many:users_boards;"`
	Members         []PublicUserInfo `json:"members" gorm:"-"`
	InvitedMembers  []PublicUserInfo `json:"invited_members" gorm:"-"`
	CardLists       []CardList       `json:"card_lists" gorm:"foreignKey:BID;constraint:OnDelete:CASCADE;"`
	Cards           []Card           `json:"-" gorm:"foreignKey:BID;constraint:OnDelete:CASCADE;"`
	Tags            []Tag            `json:"tags" gorm:"foreignKey:BID;constraint:OnDelete:CASCADE;"`
	AvailableColors []Color          `json:"colors" gorm:"-"`
}
