package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func AuthRoutes(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", controllers.Register)
		authRouter.POST("/login", controllers.Login)
		authRouter.POST("/send-otp", controllers.SendOtp)
		authRouter.POST("/verify-otp", controllers.VerifyOtp)
	}
}
