package database

import (
	"fmt"
	"time"

	"github.com/lakshya1goel/expense_tracker/models"
	_ "github.com/lib/pq"
)

func DeleteOTP(email string) error {
    return Db.Where("email = ?", email).Delete(&models.Otp{}).Error
}

func CleanExpiredOTPs() error {
    return Db.Where("otp_exp < ?", time.Now().Unix()).Delete(&models.Otp{}).Error
}

func CleanUnverifiedUsers() error {
	return Db.Where("is_verified = false").Delete(&models.User{}).Error
}

func InitDB() error {
	if Db == nil {
		return fmt.Errorf("database connection not established. Call ConnectDb first")
	}
	err := Db.AutoMigrate(&models.Expense{}, &models.User{}, &models.Otp{})
	if err != nil {
		return fmt.Errorf("error creating expenses table: %v", err)
	}

	go func() {	
		for {
			CleanExpiredOTPs()
			time.Sleep(time.Hour)
		}
	}()

	go func() {
		for {
			CleanUnverifiedUsers()
			time.Sleep(time.Hour * 24)
		}
	}()

	fmt.Println("Database tables initialized successfully")
	return nil
}
