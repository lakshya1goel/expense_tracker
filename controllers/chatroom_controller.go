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
		group := models.ChatRoom{
			Type: 	     "group",
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
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Group created successfully"})
}

func GetAllChatrooms(c *gin.Context) {
	userId := c.Param("id")
	var groups []models.ChatRoom

	if err := database.Db.Preload("Members").Where("members.id = ?", userId).Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(groups) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No groups found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": groups})
}

func GetChatroom(c *gin.Context) {
	groupId := c.Param("id")
	var group models.ChatRoom

	if err := database.Db.Preload("Members").Where("id = ?", groupId).First(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": group})
}

func AddMembers(c *gin.Context) {
	groupId := c.Param("id")
	var request dto.AddMembersDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(request.Members) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Members are required"})
		return
	}

	var group models.ChatRoom
	if err := database.Db.Where("id = ?", groupId).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if group.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cannot add members to a private chat"})
		return
	}

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
		if err := database.Db.Model(&group).Association("Members").Append(existingUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	// Handle non-existing members (e.g., send invite links)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Members added successfully"})
}

func RemoveMembers(c *gin.Context) {
	groupId := c.Param("id")
	var request dto.AddMembersDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(request.Members) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Members are required"})
		return
	}

	var group models.ChatRoom
	if err := database.Db.Where("id = ?", groupId).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if group.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cannot remove members from a private chat"})
		return
	}

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
		if err := database.Db.Model(&group).Association("Members").Delete(existingUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Members removed successfully"})
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
	var group models.ChatRoom
	if err := database.Db.Where("id = ?", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if group.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cannot update a private chat"})
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
	var group models.ChatRoom
	if err := database.Db.Where("id = ?", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if group.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cannot delete a private chat"})
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
		group := models.ChatRoom{
			Type: 	    "private",
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
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Private chat created successfully"})

}