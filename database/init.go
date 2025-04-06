package database

import (
	"fmt"

	_ "github.com/lib/pq"
)

func InitDB() error {
	if Db == nil {
		return fmt.Errorf("database connection not established. Call ConnectDb first")
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			amount INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := Db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating expenses table: %v", err)
	}

	fmt.Println("Database tables initialized successfully")
	return nil
}
