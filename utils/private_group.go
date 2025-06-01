package utils

import (
	"net/http"

	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/models"
	"gorm.io/gorm"
)

func CreatePrivateGroup(user_id uint) (int, error) {
	var user models.User
	if err := database.Db.Where("id = ?", user_id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, err
	}

	group := models.Group{
		Type:        "private",
		Name:        "Private Chat",
		Description: "",
		Users:       []*models.User{&user},
		TotalUsers:  1,
	}

	if err := database.Db.Create(&group).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	
	return http.StatusOK, nil
}
