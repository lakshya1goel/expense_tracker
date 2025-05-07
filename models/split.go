package models

import "gorm.io/gorm"

type Split struct {
	gorm.Model
	ExpenseID uint   `json:"expense_id"`
	SplitAmt   int    `json:"split_amount"`
	OwedByID  uint   `json:"owed_by_id"`
	OwedToID  uint   `json:"owed_to_id"`
	IsPaid	  bool   `json:"is_paid"`
}