package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
	"github.com/lakshya1goel/expense_tracker/middlewares"
)

func SplitRoutes(router *gin.RouterGroup) {
	splitRouter := router.Group("/split")
	splitRouter.Use(middlewares.AuthMiddleware())
	{
		splitRouter.POST("/", controllers.CreateSplit)
		splitRouter.GET("/expenses/:id", controllers.GetAllExpenses)
		splitRouter.GET("/summary/:id", controllers.GetGroupSummary)
		splitRouter.GET("/:id", controllers.GetSplit)
		splitRouter.POST("/mark-as-paid/:id", controllers.MarkSplitAsPaid)
		splitRouter.POST("/monthly-expenses/", controllers.GetMonthlyExpenses)
	}
}
