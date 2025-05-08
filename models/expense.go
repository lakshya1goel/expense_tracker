package models

import "gorm.io/gorm"

type Expense struct {
	gorm.Model
	UserID      uint     `json:"user_id"`
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Amount      int      `json:"amount"`
	GroupID     *uint    `json:"group_id"`
	PaidByCount int      `json:"paid_by_count"`
	Splits      []*Split `json:"splits"`
}
