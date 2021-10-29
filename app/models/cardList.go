package models

type CardList struct {
	CLID  uint   `json:"clid" gorm:"primaryKey"`
	BID   uint   `json:"bid" gorm:"not null;index"`
	CID   uint   `json:"cid" gorm:"foreignKey:CLID;"`
	Cards []Card `json:"cards" faker:"-" gorm:"foreignKey:CLID;"`
	Title string `json:"cardList_name" faker:"word" gorm:"not null;index"`
}
