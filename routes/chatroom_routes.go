package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func GroupRoutes(router *gin.RouterGroup) {
	authRouter := router.Group("/chatroom") 
	{
		authRouter.POST("/create-group", controllers.CreateGroup)
		authRouter.GET("/get-all/:id", controllers.GetAllChatrooms)
		authRouter.GET("/get/:id", controllers.GetChatroom)
		authRouter.POST("/add-member/:id", controllers.AddMembers)
		authRouter.POST("/remove-member/:id", controllers.RemoveMembers)
		authRouter.PUT("/update/:id", controllers.UpdateGroup)
		authRouter.DELETE("/delete/:id", controllers.DeleteGroup)
		authRouter.POST("create-private", controllers.CreatePrivateChat)
	}
}
