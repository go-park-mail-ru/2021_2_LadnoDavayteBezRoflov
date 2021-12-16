package models

type Card struct {
	CID                uint             `json:"cid" gorm:"primaryKey"`
	BID                uint             `json:"bid" gorm:"not null;index"`
	CLID               uint             `json:"clid" gorm:"not null;index"`
	PositionOnCardList uint             `json:"pos" faker:"-"`
	Title              string           `json:"card_name" faker:"word" gorm:"not null;index"`
	Description        string           `json:"description" faker:"sentence"`
	AccessPath         string           `json:"access_path"`
	DeadlineChecked    bool             `json:"deadline_check" faker:"-"`
	Deadline           string           `json:"deadline" faker:"timestamp"`
	Comments           []Comment        `json:"comments" gorm:"foreignKey:CID;constraint:OnDelete:CASCADE;"`
	CheckLists         []CheckList      `json:"check_lists" gorm:"foreignKey:CID;constraint:OnDelete:CASCADE;"`
	Tags               []Tag            `json:"tags" faker:"-" gorm:"many2many:tags_cards;"`
	Users              []User           `json:"-" gorm:"many2many:users_cards;constraint:OnDelete:CASCADE;"`
	Assignees          []PublicUserInfo `json:"assignees" gorm:"-"`
	Attachments        []Attachment     `json:"attachments" gorm:"foreignKey:CID;constraint:OnDelete:CASCADE;"`
}
