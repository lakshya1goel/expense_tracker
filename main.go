package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/routes"
)

func main() {
	database.ConnectDb()

	router := gin.Default()
	apiRouter := router.Group("/api")
	{
		routes.ExpenseRoutes(apiRouter)
		routes.AuthRoutes(apiRouter)
	}

	router.Run(":8000")
}

//TODO:
//1. resend otp rate limit
//2. resend otp after 5 minutes
//3. delete unverified user after 24 hours
