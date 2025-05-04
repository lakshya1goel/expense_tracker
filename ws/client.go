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
	User *models.User
	Room *models.ChatRoom
	mu   sync.Mutex
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		m := models.Message{
			Type: msgType, 
			Body: string(msg),
			Sender: c.User.ID,
			Room: c.Room.ID,
		}

		c.Pool.Broadcast <- m

		fmt.Println("msg recieved===>>>\n", m)
	}
}
