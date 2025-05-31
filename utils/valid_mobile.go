package utils

import (
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/models"
)

func IsValidMobile(mobile string) (bool, error) {
	result := database.Db.Where("mobile = ?", mobile).First(&models.User{})
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
