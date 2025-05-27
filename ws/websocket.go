package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var GlobalPool *Pool

func InitPool() {
	GlobalPool = NewPool()
	go GlobalPool.Start()
}

func HandleWebSocket(c *gin.Context) {
	userIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId, ok := userIdRaw.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID in context"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	client := &Client{
		Conn:   conn,
		Pool:   GlobalPool,
		UserId: uint(userId),
	}

	GlobalPool.Register <- client
	go client.Read()
}
