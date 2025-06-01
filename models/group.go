package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Type        string     `json:"type" gorm:"default:'group'"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []*User    `json:"users" gorm:"many2many:group_users"`
	Messages    []*Message `json:"messages"`
	Expenses    []*Expense `json:"expenses"`
	TotalUsers  int        `json:"total_users"`
}
