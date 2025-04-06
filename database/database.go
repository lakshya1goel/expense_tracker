package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDb() {
	var err error
	dbInfo := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", "postgres", "postgres", "expense_tracker")
	Db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic("failed to ping database: " + err.Error())
	}

	InitDB()
	fmt.Println("Connected to database")
}
