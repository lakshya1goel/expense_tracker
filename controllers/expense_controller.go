package controllers

import (
	"log"
	"net/http"

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

	err := database.Db.QueryRow(
		query,
		request.Title,
		request.Description,
		request.Amount,
	).Scan(
		&response.ID,
		&response.Title,
		&response.Description,
		&response.Amount,
		&response.CreatedAt,
	)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func GetExpenses(c *gin.Context) {
	query := `
		SELECT id, title, description, amount, created_at
		FROM expenses
		ORDER BY created_at DESC
	`

	var expenses []dto.ExpenseResponseDto

	rows, err := database.Db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roews"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var expense dto.ExpenseResponseDto
		err := rows.Scan(
			&expense.ID,
			&expense.Title,
			&expense.Description,
			&expense.Amount,
			&expense.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan expense"})
			return
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func DeleteExpense(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT EXISTS(SELECT 1 FROM expenses WHERE id = $1)
	`

	var exists bool
	err := database.Db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if expense exists"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	query = `
		DELETE FROM expenses
		WHERE id = $1
	`

	_, err = database.Db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})

}

func UpdateExpense(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT EXISTS(SELECT 1 FROM expenses WHERE id = $1)
	`

	var exists bool
	err := database.Db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if expense exists"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var request dto.CreateExpenseDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query = `
		UPDATE expenses 
		SET title = $1, description = $2, amount = $3
		WHERE id = $4
		RETURNING id, title, description, amount, created_at
	`

	var response dto.ExpenseResponseDto

	err = database.Db.QueryRow(
		query,
		request.Title,
		request.Description,
		request.Amount,
		id,
	).Scan(
		&response.ID,
		&response.Title,
		&response.Description,
		&response.Amount,
		&response.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	c.JSON(http.StatusOK, response)
}
