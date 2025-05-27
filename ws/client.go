package ws

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lakshya1goel/expense_tracker/models"
)

type Client struct {
	Conn    *websocket.Conn
	Pool    *Pool
	UserId  uint
	GroupID uint
	mu      sync.Mutex
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg models.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		switch msg.Type {
		case models.JoinGroup:
			if msg.GroupID == 0 {
				fmt.Println("Error: Invalid Group ID received.")
				continue
			}
			c.GroupID = msg.GroupID
			c.Pool.Register <- c

		case models.LeaveGroup:
			if c.GroupID == 0 {
				fmt.Println("Client has not joined any group")
				continue
			}
			c.Pool.Unregister <- c
			continue

		case models.ChatMessage, models.SplitMessage:
			if c.GroupID == 0 {
				fmt.Println("Client has not joined any group")
				continue
			}
			msg.Sender = c.UserId
			msg.GroupID = c.GroupID
			c.Pool.Broadcast <- msg

		default:
			fmt.Println("Unknown message type:", msg.Type)
		}
	}
}
