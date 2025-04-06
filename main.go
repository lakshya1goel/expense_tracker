package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/routes"
)

func main() {
	database.ConnectDb()
	log.Fatal(database.InitDB())
	
	router := gin.Default()
	api := router.Group("/api")
	{
		routes.ExpenseRoutes(api)
	}

	log.Fatal(router.Run(":8000"))
}
