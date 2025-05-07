package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []*User    `json:"users" gorm:"many2many:group_users"`
	Messages    []*Message `json:"messages"`
	TotalUsers  int         
}
