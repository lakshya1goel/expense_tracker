package dto

type CreateExpenseDto struct {
	Title       string `json:"title"`
	Description *string `json:"description"`
	Amount      int    `json:"amount"`
}

type ExpenseResponseDto struct {
	ID          uint    `json:"id"`
	Title       string `json:"title"`
	Description *string `json:"description"`
	Amount      int    `json:"amount"`
	CreatedAt   string `json:"created_at"`
}
