package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func GroupRoutes(router *gin.RouterGroup) {
	authRouter := router.Group("/group") 
	{
		authRouter.POST("/create", controllers.CreateGroup)
		authRouter.GET("/get-all/:id", controllers.GetAllGroups)
		authRouter.GET("/get/:id", controllers.GetGroup)
		// authRouter.PUT("/update/:id", controllers.UpdateGroup)
		// authRouter.DELETE("/delete/:id", controllers.DeleteGroup)
		// authRouter.POST("/add-member/:id", controllers.AddMember)
		// authRouter.POST("/remove-member/:id", controllers.RemoveMember)
	}
}
