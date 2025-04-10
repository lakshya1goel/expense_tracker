package database

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDb() {
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Invalid port number: " + err.Error())
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	configData := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	var dbErr error
	Db, dbErr = gorm.Open(postgres.Open(configData), &gorm.Config{})
	if dbErr != nil {
		panic("Error connecting to database: " + dbErr.Error())
	}

	fmt.Println("\x1b[32m...............Database connected..................\x1b[0m")

	InitDB()
}
