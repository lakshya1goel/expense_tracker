package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"gorm.io/gorm"
)

func CreateGroup(c *gin.Context) {
	var request dto.CreateChatDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Name is required"})
		return
	}

	//TODO: if Users is empty then return the response with invite link to group

	var existingUsers []*models.User
	var nonExistingUsers []string

	for _, member := range request.Users {
		if member == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Member phone no. is required"})
			return
		}

		var user models.User
		result := database.Db.Where("mobile = ?", member).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				nonExistingUsers = append(nonExistingUsers, member)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
				return
			}
		} else {
			existingUsers = append(existingUsers, &user)
		}
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var creator models.User
	result := database.Db.Where("id = ?", userId).First(&creator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
		return
	}

	existingUsers = append(existingUsers, &creator)

	if len(existingUsers) > 0 {
		group := models.Group{
			Name:        request.Name,
			Description: request.Description,
			Users:       existingUsers,
			TotalUsers:  len(existingUsers),
		}
		if err := database.Db.Create(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	//TODO: Handle non-existing Users (e.g., send invite links)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Group created successfully"})
}

func GetAllGroups(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var groupIDs []uint
	if err := database.Db.Table("group_users").Where("user_id = ?", userId).Pluck("group_id", &groupIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(groupIDs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No groups found"})
		return
	}

	var groups []models.Group
	if err := database.Db.Where("id IN ?", groupIDs).Preload("Users").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var response []gin.H
	for _, group := range groups {
		response = append(response, gin.H{
			"id":          group.ID,
			"name":        group.Name,
			"description": group.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Groups fetched successfully!", "data": response})
}

func GetGroupHistory(c *gin.Context) {
	groupId := c.Param("id")
	var group models.Group
	if err := database.Db.Preload("Expenses").Preload("Messages").Where("id = ?", groupId).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"expenses":    group.Expenses,
		"messages":    group.Messages,
		"total_users": group.TotalUsers,
	}})
}

func GetGroup(c *gin.Context) {
	groupId := c.Param("id")
	var group models.Group

	if err := database.Db.Preload("Users").Where("id = ?", groupId).First(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": group})
}

func AddUsers(c *gin.Context) {
	groupId := c.Param("id")
	var request dto.AddUsersDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(request.Users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Users are required"})
		return
	}

	var group models.Group
	if err := database.Db.Where("id = ?", groupId).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var existingUsers []models.User
	var nonExistingUsers []string

	for _, member := range request.Users {
		if member == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Member phone no. is required"})
			return
		}

		var user models.User
		result := database.Db.Where("phone = ?", member).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				nonExistingUsers = append(nonExistingUsers, member)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
				return
			}
		} else {
			existingUsers = append(existingUsers, user)
		}
	}

	if len(existingUsers) > 0 {
		if err := database.Db.Model(&group).Association("Users").Append(existingUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	// Handle non-existing Users (e.g., send invite links)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Users added successfully"})
}

func RemoveUsers(c *gin.Context) {
	groupId := c.Param("id")
	var request dto.AddUsersDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(request.Users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Users are required"})
		return
	}

	var group models.Group
	if err := database.Db.Where("id = ?", groupId).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var existingUsers []models.User
	var nonExistingUsers []string

	for _, member := range request.Users {
		if member == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Member phone no. is required"})
			return
		}

		var user models.User
		result := database.Db.Where("phone = ?", member).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				nonExistingUsers = append(nonExistingUsers, member)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
				return
			}
		} else {
			existingUsers = append(existingUsers, user)
		}
	}

	if len(existingUsers) > 0 {
		if err := database.Db.Model(&group).Association("Users").Delete(existingUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Users removed successfully"})
}

func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var request dto.UpdateGroupDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Name is required"})
		return
	}
	var group models.Group
	if err := database.Db.Where("id = ?", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	group.Name = request.Name
	group.Description = request.Description
	group.UpdatedAt = time.Now()
	if err := database.Db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Group updated successfully"})
}

func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group
	if err := database.Db.Where("id = ?", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := database.Db.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Group deleted successfully"})
}

func CreatePrivateChat(c *gin.Context) {
	var request dto.CreateChatDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var existingUsers []*models.User
	var nonExistingUsers []string

	for _, member := range request.Users {
		if member == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Member phone no. is required"})
			return
		}

		var user models.User
		result := database.Db.Where("phone = ?", member).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				nonExistingUsers = append(nonExistingUsers, member)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
				return
			}
		} else {
			existingUsers = append(existingUsers, &user)
		}
	}

	if len(existingUsers) > 0 {
		group := models.Group{
			Name:        request.Name,
			Description: request.Description,
			Users:       existingUsers,
		}
		if err := database.Db.Create(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	// Handle non-existing Users (e.g., send invite links)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Private chat created successfully"})

}
