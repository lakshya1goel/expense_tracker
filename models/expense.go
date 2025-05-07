package models

import "gorm.io/gorm"

type Expense struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	GroupID     *uint
	UserID      uint
}
