package main

import (
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

	router.Run(":8000")
}
