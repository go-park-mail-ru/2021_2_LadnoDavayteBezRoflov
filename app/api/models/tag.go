package models

type Tag struct {
	TGID    uint   `json:"tgid" gorm:"primaryKey"`
	BID     uint   `json:"bid" gorm:"not null;index"`
	Title   string `json:"tag_name" faker:"word" gorm:"not null;index"`
	ColorID uint   `json:"-" gorm:"not null;index"`
	Color   Color  `json:"color" gorm:"-"`
	Cards   []Card `json:"-" faker:"-" gorm:"many2many:tags_cards;"`
}
