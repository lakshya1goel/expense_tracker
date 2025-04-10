package database

import (
	"fmt"

	"github.com/lakshya1goel/expense_tracker/models"
	_ "github.com/lib/pq"
)

func InitDB() error {
	if Db == nil {
		return fmt.Errorf("database connection not established. Call ConnectDb first")
	}
	err := Db.AutoMigrate(&models.Expense{}, &models.User{}, &models.Otp{})
	if err != nil {
		return fmt.Errorf("error creating expenses table: %v", err)
	}

	fmt.Println("Database tables initialized successfully")
	return nil
}
