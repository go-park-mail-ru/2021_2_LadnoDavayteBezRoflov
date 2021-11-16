package models

type Card struct {
	CID                uint   `json:"cid" gorm:"primaryKey"`
	BID                uint   `json:"bid" gorm:"not null;index"`
	CLID               uint   `json:"clid" gorm:"not null;index"`
	PositionOnCardList uint   `json:"pos" faker:"-"`
	Title              string `json:"card_name" faker:"word" gorm:"not null;index"`
	Description        string `json:"description" faker:"sentence"`
	// Deadline           time.Time `json:"-" faker:"-"`
	DeadlineChecked bool   `json:"deadline_check" faker:"-"`
	Deadline        string `json:"deadline" faker:"timestamp"`
}
