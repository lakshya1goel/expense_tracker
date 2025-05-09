package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
	"github.com/lakshya1goel/expense_tracker/middlewares"
	"github.com/lakshya1goel/expense_tracker/ws"
)

func GroupRoutes(router *gin.RouterGroup) {
	groupRouter := router.Group("/group")
	groupRouter.Use(middlewares.AuthMiddleware())
	{
		groupRouter.GET("/ws", ws.HandleWebSocket)
		groupRouter.POST("/", controllers.CreateGroup)
		groupRouter.GET("/get-all/:id", controllers.GetAllGroups)
		groupRouter.GET("/get/:id", controllers.GetGroup)
		groupRouter.POST("/add-member/:id", controllers.AddUsers)
		groupRouter.DELETE("/remove-member/:id", controllers.RemoveUsers)
		groupRouter.PUT("/:id", controllers.UpdateGroup)
		groupRouter.DELETE("/:id", controllers.DeleteGroup)
		groupRouter.POST("create-private", controllers.CreatePrivateChat)
	}
}
