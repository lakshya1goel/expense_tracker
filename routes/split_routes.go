package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func SplitRoutes(router *gin.RouterGroup) {
	splitRouter := router.Group("/split")
	{
		splitRouter.POST("/", controllers.CreateSplit)
	}
}
