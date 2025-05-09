package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email            string     `json:"email" gorm:"unique"`
	Mobile           string     `json:"mobile" gorm:"unique"`
	Password         string     `json:"password"`
	IsEmailVerified  bool       `json:"is_email_verified" gorm:"default:false"`
	IsMobileVerified bool       `json:"is_mobile_verified" gorm:"default:false"`
	Groups           []*Group   `json:"groups" gorm:"many2many:group_users"`
	Expenses         []*Expense `json:"expenses"`
}
