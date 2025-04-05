package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

)

func ConnectDb() {
	dbInfo := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", "postgres", "postgres", "expense_tracker")
	fmt.Println(dbInfo)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic("failed to ping database: " + err.Error())
	}
	
	fmt.Println("Connected to database")
}
