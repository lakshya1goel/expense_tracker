package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/middlewares"
	"github.com/lakshya1goel/expense_tracker/ws"
)

func WsRoutes(router *gin.RouterGroup) {
	wsRouter := router.Group("/ws")
	wsRouter.Use(middlewares.AuthMiddleware())
	{
		wsRouter.GET("/", ws.HandleWebSocket)
	}
}
