package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/routes"
)

func main() {
	database.ConnectDb()

	// if err := database.InitDB(); err != nil {
	// 	log.Fatal(err)
	// }

	router := gin.Default()
	api := router.Group("/api")
	{
		routes.ExpenseRoutes(api)
	}

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
