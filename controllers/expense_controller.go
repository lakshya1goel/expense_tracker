package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"gorm.io/gorm"
)

func CreateExpense(c *gin.Context) {
	var request dto.CreateExpenseDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	expense := models.Expense{
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
	}

	result := database.Db.Create(&expense)

	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create expense " + result.Error.Error()})
		return
	}

	response := dto.ExpenseResponseDto{
		ID:        expense.ID,
		Title:     expense.Title,
		Amount:    expense.Amount,
		CreatedAt: expense.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": response})
}

func GetExpenses(c *gin.Context) {
	var dbExpenses []models.Expense
	result := database.Db.Order("created_at").Find(&dbExpenses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch expenses " + result.Error.Error()})
		return
	}

	expenses := make([]dto.ExpenseResponseDto, len(dbExpenses))
	for i, expense := range dbExpenses {
		expenses[i] = dto.ExpenseResponseDto{
			ID:          expense.ID,
			Title:       expense.Title,
			Description: expense.Description,
			Amount:      expense.Amount,
			CreatedAt:   expense.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": expenses})
}

func DeleteExpense(c *gin.Context) {
	id := c.Param("id")

	expenseResult := database.Db.First(&models.Expense{}, id)

	if expenseResult.Error != nil {
		if errors.Is(expenseResult.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Expense not found " + expenseResult.Error.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB Error " + expenseResult.Error.Error()})
		}
		return
	}

	result := database.Db.Delete(&models.Expense{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete expense " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Expense deleted successfully"})
}

func UpdateExpense(c *gin.Context) {
	id := c.Param("id")

	var request dto.CreateExpenseDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	expense := models.Expense{
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
	}

	var existingExpense models.Expense
	result := database.Db.First(&existingExpense, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Expense not found " + result.Error.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "DB Error " + result.Error.Error()})
		}
		return
	}

	existingExpense.Title = expense.Title
	existingExpense.Description = expense.Description
	existingExpense.Amount = expense.Amount
	database.Db.Save(&existingExpense)

	response := dto.ExpenseResponseDto{
		ID:        existingExpense.ID,
		Title:     existingExpense.Title,
		Amount:    existingExpense.Amount,
		CreatedAt: existingExpense.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}
