package models

type User struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Email         string `json:"email" gorm:"unique"`
	Password      string `json:"password"`
	IsVerified    bool   `json:"is_verified" gorm:"default:false"`
}