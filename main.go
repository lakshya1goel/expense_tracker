package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/routes"
	"github.com/lakshya1goel/expense_tracker/utils"
	"github.com/lakshya1goel/expense_tracker/ws"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	database.ConnectDb()
	utils.InitGoogleOAuth()
	ws.InitPool()

	router := gin.Default()
	apiRouter := router.Group("/api")
	{
		routes.ExpenseRoutes(apiRouter)
		routes.AuthRoutes(apiRouter)
		routes.OauthRoutes(apiRouter)
		routes.GroupRoutes(apiRouter)
		routes.SplitRoutes(apiRouter)
	}

	router.Run(":8000")
}

//TODO:
//1. resend otp rate limit
//2. resend otp after 5 minutes
