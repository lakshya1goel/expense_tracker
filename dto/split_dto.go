package dto

import (
	"time"

	"github.com/lakshya1goel/expense_tracker/models"
)

type CreateSplitDto struct {
	GroupID     uint    `json:"group_id"`
	Amount      int     `json:"amount"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type CreatePersonalExpenseDto struct {
	GroupID     uint    `json:"group_id"`
	Amount      int     `json:"amount"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type SplitResponseDto struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Description *string         `json:"description"`
	Amount      int             `json:"amount"`
	GroupID     uint            `json:"group_id"`
	PaidByCount int             `json:"PaidByCount"`
	Splits      []*models.Split `json:"splits"`
}

type SplitWsDto struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	ExpenseAmount int       `json:"expense_amount"`
	SplitAmt      int       `json:"split_amount"`
	SenderID      uint      `json:"sender_id"`
	OwedByIDs     []uint    `json:"owed_by_ids"`
}

type MarkSplitAsPaidDto struct {
	ExpenseID uint `json:"expense_id"`
	OwedByID  uint `json:"owed_by_id"`
	OwedToID  uint `json:"owed_to_id"`
}

type MonthlyExpenseRequestDto struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type MonthlyExpenseResponseDto struct {
	Month        int     `json:"month"`
	Year         int     `json:"year"`
	SpentAmout   float64 `json:"spent_amount"`
	OwedToAmount float64 `json:"owed_to_amount"`
	OwedByAmount float64 `json:"owed_by_amount"`
}
