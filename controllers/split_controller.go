package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"github.com/lakshya1goel/expense_tracker/ws"
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

	var userIDs []uint
	if err := database.Db.Table("group_users").Where("group_id = ?", request.GroupID).Pluck("user_id", &userIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	expense := models.Expense{
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
		GroupID:     &request.GroupID,
		UserID:      uint(userID.(float64)),
		PaidByCount: 1,
	}
	if err := database.Db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	splitAmt := expense.Amount / len(userIDs)
	for _, userId := range userIDs {
		var split models.Split
		if userId == uint(userID.(float64)) {
			split = models.Split{
				ExpenseID: expense.ID,
				GroupID:   *expense.GroupID,
				SplitAmt:  splitAmt,
				OwedByID:  userId,
				OwedToID:  uint(userID.(float64)),
				IsPaid:    true,
			}
		} else {
			split = models.Split{
				ExpenseID: expense.ID,
				GroupID:   *expense.GroupID,
				SplitAmt:  splitAmt,
				OwedByID:  userId,
				OwedToID:  uint(userID.(float64)),
				IsPaid:    false,
			}
		}
		if err := database.Db.Create(&split).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	response := dto.SplitResponseDto{
		ID:          expense.ID,
		Title:       expense.Title,
		Description: expense.Description,
		Amount:      expense.Amount,
		GroupID:     *expense.GroupID,
		PaidByCount: expense.PaidByCount,
		Splits:      expense.Splits,
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Split created successfully", "data": response})

	var owedByIds []uint
	for _, userId := range userIDs {
		if userId != uint(userID.(float64)) {
			owedByIds = append(owedByIds, userId)
		}
	}
	splitWsDto := dto.SplitWsDto{
		ID:            expense.ID,
		CreatedAt:     expense.CreatedAt,
		ExpenseAmount: expense.Amount,
		SplitAmt:      splitAmt,
		SenderID:      uint(userID.(float64)),
		OwedByIDs:     owedByIds,
	}

	splitBodyBytes, err := json.Marshal(splitWsDto)
	if err != nil {
		fmt.Printf("Error marshaling splitWsDto:", err)
		return
	}

	ws.GlobalPool.Broadcast <- models.Message{
		Type:    models.SplitMessage,
		Body:    string(splitBodyBytes),
		GroupID: request.GroupID,
		Sender:  uint(userID.(float64)),
	}
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
	var request dto.MarkSplitAsPaidDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	var split models.Split
	if err := database.Db.Where("expense_id = ? AND owed_to_id = ? AND owed_by_id = ?", request.ExpenseID, request.OwedToID, request.OwedByID).First(&split).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Split not found"})
			return
		}
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

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var group models.Group
	if err := database.Db.Preload("Users").Where("id = ?", groupId).First(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var settlements [](map[string]interface{})
	var totalGroupOwedTo, totalGroupOwedBy float64

	for _, user := range group.Users {
		if user.ID == uint(userId.(float64)) {
			continue
		}

		var owedToAmt, owedByAmt float64
		if err := database.Db.Model(&models.Split{}).Where("owed_to_id = ? AND owed_by_id = ? AND group_id = ? AND is_paid = false", uint(userId.(float64)), user.ID, group.ID).Select("COALESCE(SUM(split_amt), 0)").Scan(&owedToAmt).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := database.Db.Model(&models.Split{}).Where("owed_by_id = ? AND owed_to_id = ? AND group_id = ? AND is_paid = false", uint(userId.(float64)), user.ID, group.ID).Select("COALESCE(SUM(split_amt), 0)").Scan(&owedByAmt).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		var noOfSplits int64
		if err := database.Db.Model(&models.Split{}).Where("(owed_by_id = ? OR owed_to_id = ?) AND NOT (owed_by_id = ? AND owed_to_id = ?)", uint(userId.(float64)), uint(userId.(float64)), uint(userId.(float64)), uint(userId.(float64))).Count(&noOfSplits).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		if owedToAmt > owedByAmt {
			totalGroupOwedTo += owedToAmt - owedByAmt
		} else {
			totalGroupOwedBy += owedByAmt - owedToAmt
		}

		settlement := map[string]interface{}{
			"user_id":      user.ID,
			"user_name":    user.Email,
			"owed_to_amt":  owedToAmt,
			"owed_by_amt":  owedByAmt,
			"settlement":   owedToAmt - owedByAmt,
			"no_of_splits": noOfSplits,
		}

		settlements = append(settlements, settlement)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group summary fetched successfully",
		"data": map[string]interface{}{
			"settlements":         settlements,
			"total_group_owed_to": totalGroupOwedTo,
			"total_group_owed_by": totalGroupOwedBy,
		},
	})
}

func GetMonthlyExpenses(c *gin.Context) {
	var request dto.MonthlyExpenseRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var spentAmt float64
	var owedByAmt float64
	var owedToAmt float64
	if err := database.Db.Model(&models.Split{}).Where("owed_by_id = ? AND owed_to_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", uint(userId.(float64)), uint(userId.(float64)), request.Month, request.Year).Select("COALESCE(SUM(split_amt), 0)").Scan(&spentAmt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := database.Db.Model(&models.Split{}).Where("owed_to_id = ? AND owed_by_id != ? AND is_paid = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", uint(userId.(float64)), uint(userId.(float64)), false, request.Month, request.Year).Select("COALESCE(SUM(split_amt), 0)").Scan(&owedToAmt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := database.Db.Model(&models.Split{}).Where("owed_by_id = ? AND owed_to_id != ? AND is_paid = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", uint(userId.(float64)), uint(userId.(float64)), false, request.Month, request.Year).Select("COALESCE(SUM(split_amt), 0)").Scan(&owedByAmt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	response := dto.MonthlyExpenseResponseDto{
		Month:        request.Month,
		Year:         request.Year,
		SpentAmout:   spentAmt,
		OwedToAmount: owedToAmt,
		OwedByAmount: owedByAmt,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Monthly expenses fetched successfully", "data": response})
}
