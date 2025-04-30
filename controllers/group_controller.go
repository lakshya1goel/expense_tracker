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
	var request dto.CreateGroupDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Name is required"})
		return
	}

	//TODO: if members is empty then return the response with invite link to group

	var existingUsers []models.User
	var nonExistingMembers []string

	for _, member := range request.Members {
		if member == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Member phone no. is required"})
			return
		}

		var user models.User
		result := database.Db.Where("phone = ?", member).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				nonExistingMembers = append(nonExistingMembers, member)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": result.Error.Error()})
				return
			}
		} else {
			existingUsers = append(existingUsers, user)
		}
	}

	if len(existingUsers) > 0 {
		group := models.Group{
			Name:        request.Name,
			Description: request.Description,
			Members:     existingUsers,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := database.Db.Create(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	// Handle non-existing members (e.g., send invite links)
	// for _, member := range nonExistingMembers {
	// 	// TODO: send invite link to the user
	// }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Group created successfully"})
}

func GetAllGroups(c *gin.Context) {
	userId := c.Param("id")
	var groups []models.Group

	if err := database.Db.Preload("Members").Where("members.id = ?", userId).Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(groups) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No groups found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "groups": groups})
}

func GetGroup(c *gin.Context) {
	groupId := c.Param("id")
	var group models.Group

	if err := database.Db.Preload("Members").Where("id = ?", groupId).First(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "group": group})
}