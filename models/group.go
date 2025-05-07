package models

import "gorm.io/gorm"

// TODO: reanme it to Group
type Group struct {
	gorm.Model
	Type        string     `json:"type" gorm:"type:ENUM('group', 'private')"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []*User    `json:"users" gorm:"many2many:group_users"`
	Messages    []*Message `json:"messages"`
	TotalUsers  int         
}
