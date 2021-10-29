package models

import "gorm.io/datatypes"

type Card struct {
	CID         uint           `json:"cid" gorm:"primaryKey"`
	BID         uint           `json:"bid" gorm:"not null;index"`
	CLID        uint           `json:"clid" gorm:"not null;index"`
	Title       string         `json:"card_name" faker:"word" gorm:"not null;index"`
	Description string         `json:"description" faker:"sentence"`
	Deadline    datatypes.Date `json:"deadline" faker:"timestamp"`
}
