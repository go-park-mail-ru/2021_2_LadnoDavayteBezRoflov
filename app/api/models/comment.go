package models

import "time"

type Comment struct {
	CMID       uint           `json:"cmid" gorm:"primaryKey"`
	CID        uint           `json:"cid" gorm:"not null;index"`
	UID        uint           `json:"uid" gorm:"not null;index"`
	Text       string         `json:"text" faker:"sentence" gorm:"not null;index"`
	Date       time.Time      `json:"-" faker:"-"`
	DateParsed string         `json:"date" faker:"-" gorm:"-"`
	User       PublicUserInfo `json:"user" gorm:"-"`
}
