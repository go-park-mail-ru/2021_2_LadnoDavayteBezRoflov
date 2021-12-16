package models

type Attachment struct {
	ATID               uint   `json:"atid" gorm:"primaryKey"`
	CID                uint   `json:"cid" gorm:"not null;index"`
	AttachmentTechName string `json:"file_tech_name" gorm:"not null;index"`
	AttachmentPubName  string `json:"file_pub_name" gorm:"not null;index"`
}
