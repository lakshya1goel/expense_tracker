package models

import "time"

type Expense struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}
