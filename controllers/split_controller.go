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
		PaidByCount:      1,
	}
	if err := database.Db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	splitAmt := expense.Amount / len(group.Users)
	for _, user := range group.Users {
		var split models.Split
		if user.ID == userID.(uint) {
			split = models.Split{
				ExpenseID: expense.ID,
				GroupID:   *expense.GroupID,
				SplitAmt:  splitAmt,
				OwedByID:  user.ID,
				OwedToID:  userID.(uint),
				IsPaid:    true,
			}
		}
		split = models.Split{
			ExpenseID: expense.ID,
			GroupID:   *expense.GroupID,
			SplitAmt:  splitAmt,
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
		PaidByCount:  expense.PaidByCount,
		Splits:      expense.Splits,
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Split created successfully", "data": response})
}

func GetAllExpenses(c *gin.Context) {
	groupId := c.Param("id")
	var expenses []models.Expense
	if err := database.Db.Where("group_id = ?", groupId).Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(expenses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No expenses found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Expenses fetched successfully", "data": expenses})
}

func GetSplit(c *gin.Context) {
	expenseId := c.Param("id")
	var splits []models.Split
	if err := database.Db.Where("expense_id = ?", expenseId).Find(&splits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(splits) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No splits found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Splits fetched successfully", "data": splits})
}

func MarkSplitAsPaid(c *gin.Context) {
	splitId := c.Param("id")
	var split models.Split
	if err := database.Db.Where("id = ?", splitId).First(&split).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if split.IsPaid {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Split already paid"})
		return
	}
	split.IsPaid = true
	if err := database.Db.Save(&split).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var expense models.Expense
	if err := database.Db.Where("id = ?", split.ExpenseID).First(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := database.Db.Model(&expense).Update("paid_by_count", expense.PaidByCount+1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Split marked as paid successfully"})
}

func GetGroupSummary(c *gin.Context) {
	groupId := c.Param("id")
	var group models.Group

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var groupUsers []models.User
	if err := database.Db.Where("group_id = ?", groupId).Find(&groupUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	for _, user := range groupUsers {
		if user.ID == userId.(uint) {
			group.Users = append(group.Users, &user)
		}
	}
}