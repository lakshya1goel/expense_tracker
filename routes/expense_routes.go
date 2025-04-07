package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func ExpenseRoutes(router *gin.RouterGroup) {
	expense := router.Group("/expenses")
	{
		expense.POST("/", controllers.CreateExpense)
		expense.GET("/", controllers.GetExpenses)
		expense.DELETE("/:id", controllers.DeleteExpense)
		expense.PUT("/:id", controllers.UpdateExpense)
	}
}
