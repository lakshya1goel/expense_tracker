package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/controllers"
)

func OauthRoutes(router *gin.RouterGroup) {
	oauthRouter := router.Group("/oauth")
	{
		oauthRouter.GET("/google/login", controllers.GoogleSignIn)
		oauthRouter.GET("/google/callback", controllers.GoogleCallback)
		oauthRouter.POST("/google/idToken", controllers.VerifyIdToken)
	}
}
