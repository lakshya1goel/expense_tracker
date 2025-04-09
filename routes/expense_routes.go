package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func ExpenseRoutes(router *gin.RouterGroup) {
	expenseRouter := router.Group("/expenses")
	{
		expenseRouter.POST("/", controllers.CreateExpense)
		expenseRouter.GET("/", controllers.GetExpenses)
		expenseRouter.DELETE("/:id", controllers.DeleteExpense)
		expenseRouter.PUT("/:id", controllers.UpdateExpense)
	}
}
