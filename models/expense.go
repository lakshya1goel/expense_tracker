package models

import "gorm.io/gorm"

type Expense struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	GroupID     *uint  `json:"group_id"`
	PaidBy      uint   `json:"user_id"`
	Splits      []*Split `json:"splits"`
}
