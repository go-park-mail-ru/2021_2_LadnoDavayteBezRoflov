package models

type CheckList struct {
	CHLID          uint            `json:"chlid" gorm:"primaryKey"`
	CID            uint            `json:"cid" gorm:"not null;index"`
	Title          string          `json:"title" faker:"word" gorm:"not null;index"`
	CheckListItems []CheckListItem `json:"check_list_items" gorm:"foreignKey:CHLID;constraint:OnDelete:CASCADE;"`
}
