package ws

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/lakshya1goel/expense_tracker/models"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Rooms      map[string]map[*Client]bool
	Broadcast  chan models.Message
	mu         sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan models.Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.mu.Lock()
			roomId := strconv.Itoa(int(client.Room.ID))
			if _, ok := pool.Rooms[roomId]; !ok {
				pool.Rooms[roomId] = make(map[*Client]bool)
			}
			pool.Rooms[roomId][client] = true
			pool.mu.Unlock()

			fmt.Printf("New user joined room: %s, total users: %d\n", roomId, len(pool.Rooms[roomId]))

			for c := range pool.Rooms[roomId] {
				if c != client {
					err := c.Conn.WriteJSON(models.Message{
						Type:   1,
						Body:   "New User Joined",
						Sender: 0,
						Room:   client.Room.ID,
					})
					if err != nil {
						fmt.Println("Write error:", err)
					}
				}
			}
		case client := <-pool.Unregister:
			pool.mu.Lock()
			roomId := strconv.Itoa(int(client.Room.ID))
			if _, ok := pool.Rooms[roomId]; ok {
				delete(pool.Rooms[roomId], client)
				if len(pool.Rooms[roomId]) == 0 {
					delete(pool.Rooms, roomId)
				}
			}
			pool.mu.Unlock()

			for c := range pool.Rooms[roomId] {
				err := c.Conn.WriteJSON(models.Message{
					Type:   1,
					Body:   "User Disconnected",
					Sender: 0,
					Room:   client.Room.ID,
				})
				if err != nil {
					fmt.Println("Write error:", err)
				}
			}
		case msg := <-pool.Broadcast:
			fmt.Println("Broadcasting message to room:", msg.Room)
			pool.mu.RLock()
			roomId := strconv.Itoa(int(msg.Room))
			for c := range pool.Rooms[roomId] {
				if c.User.ID == msg.Sender {
					continue
				}
				if err := c.Conn.WriteJSON(msg); err != nil {
					fmt.Println("Broadcast write error:", err)
				}
			}
			pool.mu.RUnlock()
		}
	}
}
