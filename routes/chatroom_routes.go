package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
	"github.com/lakshya1goel/expense_tracker/ws"
)

func ChatRoutes(router *gin.RouterGroup) {
	chatRouter := router.Group("/chatroom")
	{
		chatRouter.GET("/ws", ws.HandleWebSocket)
		chatRouter.POST("/create-group", controllers.CreateGroup)
		chatRouter.GET("/get-all/:id", controllers.GetAllChatrooms)
		chatRouter.GET("/get/:id", controllers.GetChatroom)
		chatRouter.POST("/add-member/:id", controllers.AddMembers)
		chatRouter.POST("/remove-member/:id", controllers.RemoveMembers)
		chatRouter.PUT("/update/:id", controllers.UpdateGroup)
		chatRouter.DELETE("/delete/:id", controllers.DeleteGroup)
		chatRouter.POST("create-private", controllers.CreatePrivateChat)
	}
}
