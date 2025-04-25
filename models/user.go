package models

type User struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Email         string `json:"email" gorm:"unique"`
	Mobile	  string `json:"mobile_no" gorm:"unique"`
	Password      string `json:"password"`
	IsEmailVerified    bool   `json:"is_email_verified" gorm:"default:false"`
	IsMobileVerified   bool   `json:"is_mobile_verified" gorm:"default:false"`
}