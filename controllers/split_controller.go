package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"gorm.io/gorm"
)

func CreateSplit(c *gin.Context) {
	var request dto.CreateSplitDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.GroupID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Group ID is required"})
		return
	}

	var group models.Group
	result := database.Db.Where("group_id = ?", request.GroupID).First(group)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	expense := models.Expense{
		Title:       request.Title,
		Description: *request.Description,
		Amount:      request.Amount,
		GroupID:     &request.GroupID,
		PaidBy:      userID.(uint),
	}
	if err := database.Db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	splitAmt := expense.Amount / len(group.Users)
	for _, user := range group.Users {
		if user.ID == userID.(uint) {
			continue
		}
		split := models.Split{
			ExpenseID: expense.ID,
			SplitAmt: splitAmt,
			OwedByID:  user.ID,
			OwedToID:  userID.(uint),
			IsPaid:    false,
		}
		if err := database.Db.Create(&split).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	response := dto.SplitResponseDto{
		ID:          expense.ID,
		Title:       expense.Title,
		Description: &expense.Description,
		Amount:      expense.Amount,
		GroupID:     *expense.GroupID,
		PaidBy:      expense.PaidBy,
		Splits:      expense.Splits,	
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Split created successfully", "data": response})
}
