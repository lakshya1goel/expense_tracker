package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
)

func CreateExpense(c *gin.Context) {
	var request dto.CreateExpenseDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO expenses (title, description, amount)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, amount, created_at
	`

	var response dto.ExpenseResponseDto
	var createdAt time.Time

	err := database.DB.QueryRow(
		query,
		request.Title,
		request.Description,
		request.Amount,
	).Scan(
		&response.ID,
		&response.Title,
		&response.Description,
		&response.Amount,
		&createdAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	c.JSON(http.StatusCreated, response)
}
