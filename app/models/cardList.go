package models

type CardList struct {
	CLID            uint   `json:"clid" gorm:"primaryKey"`
	BID             uint   `json:"bid" gorm:"not null;index"`
	CID             uint   `json:"cid" gorm:"foreignKey:CLID;"`
	PositionOnBoard uint   `json:"pos" faker:"-"`
	Cards           []Card `json:"cards" gorm:"foreignKey:CLID;constraint:OnDelete:CASCADE;"`
	Title           string `json:"cardList_name" faker:"word" gorm:"not null;index"`
}
