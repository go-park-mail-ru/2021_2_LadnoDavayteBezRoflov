package models

type CheckListItem struct {
	CHLIID uint   `json:"chliid" gorm:"primaryKey"`
	CHLID  uint   `json:"chlid" gorm:"not null;index"`
	Text   string `json:"text" faker:"word" gorm:"not null;index"`
	Status bool   `json:"status" gorm:"not null"`
}
