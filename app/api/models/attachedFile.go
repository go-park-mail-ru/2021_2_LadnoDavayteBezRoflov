package models

type AttachedFile struct {
	ATID           uint   `json:"atid" gorm:"primaryKey"`
	CID            uint   `json:"cid" gorm:"not null;index"`
	AttachmentTech string `json:"file_tech_name" gorm:"not null;index"`
	AttachmentPub  string `json:"file_pub_name" gorm:"not null;index"`
}
