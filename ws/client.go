package ws

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lakshya1goel/expense_tracker/models"
)

type Client struct {
	Conn *websocket.Conn
	Pool *Pool
	UserId uint
	RoomId string
	mu   sync.Mutex
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
		case models.JoinRoom:
			if msg.Room == "" {
				fmt.Println("Error: Invalid Room ID received.")
				continue
			}
			c.RoomId = msg.Room
			c.Pool.Register <- c

		case models.LeaveRoom:
			if c.RoomId == "" {
				fmt.Println("Client has not joined any room")
				continue
			}
			c.Pool.Unregister <- c
			return

		case models.ChatMessage:
			if c.RoomId == "" {
				fmt.Println("Client has not joined any room")
				continue
			}
			msg.Sender = c.UserId
			msg.Room = c.RoomId
			c.Pool.Broadcast <- msg

		default:
			fmt.Println("Unknown message type:", msg.Type)
		}
	}
}
