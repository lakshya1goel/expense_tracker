package dto

import "github.com/lakshya1goel/expense_tracker/models"

type CreateSplitDto struct {
	GroupID     uint    `json:"group_id"`
	Amount      int     `json:"amount"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type SplitResponseDto struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Amount      int      `json:"amount"`
	GroupID     uint     `json:"group_id"`
	PaidByCount      int     `json:"PaidByCount"`
	Splits      []*models.Split `json:"splits"`
}
