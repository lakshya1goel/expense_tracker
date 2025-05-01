package models

import "time"

type ChatRoom struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Type        string    `json:"type" gorm:"type:ENUM('group', 'private')"`
	Name        string   `json:"name"`
	Description string    `json:"description"`
	Members     []User    `json:"members" gorm:"many2many:group_users"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
