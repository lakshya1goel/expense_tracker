package database

import (
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDb() {
	configData := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", "postgres", "postgres", "expense_tracker")
	var dbErr error
	Db, dbErr = gorm.Open(postgres.Open(configData), &gorm.Config{})
	if dbErr != nil {
		panic("Error connecting to database: " + dbErr.Error())
	}

	fmt.Println("\x1b[32m...............Database connected..................\x1b[0m")
}
