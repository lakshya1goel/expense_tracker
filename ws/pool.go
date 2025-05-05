package ws

import (
	"fmt"
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
			if client.RoomId == "" {
				fmt.Println("Error: Client has not joined a valid room. Skipping registration.")
				continue
			}

			pool.mu.Lock()
			if _, ok := pool.Rooms[client.RoomId]; !ok {
				pool.Rooms[client.RoomId] = make(map[*Client]bool)
			}
			pool.Rooms[client.RoomId][client] = true
			pool.mu.Unlock()

			fmt.Printf("New user joined room: %s, total users: %d\n", client.RoomId, len(pool.Rooms[client.RoomId]))

			for c := range pool.Rooms[client.RoomId] {
				if c != client {
					err := c.Conn.WriteJSON(models.Message{
						Type:   models.JoinRoom,
						Body:   "New User Joined",
						Sender: 0,
						Room:   client.RoomId,
					})
					if err != nil {
						fmt.Println("Write error:", err)
					}
				}
			}
		case client := <-pool.Unregister:
			if client.RoomId == "" {
				fmt.Println("Error: Client with invalid RoomId tried to unregister.")
				continue
			}

			pool.mu.Lock()
			if _, ok := pool.Rooms[client.RoomId]; ok {
				delete(pool.Rooms[client.RoomId], client)
				if len(pool.Rooms[client.RoomId]) == 0 {
					delete(pool.Rooms, client.RoomId)
				}
			}
			pool.mu.Unlock()

			for c := range pool.Rooms[client.RoomId] {
				err := c.Conn.WriteJSON(models.Message{
					Type:   models.LeaveRoom,
					Body:   "User Disconnected",
					Sender: 0,
					Room:   client.RoomId,
				})
				if err != nil {
					fmt.Println("Write error:", err)
				}
			}
		case msg := <-pool.Broadcast:
			fmt.Println("Broadcasting message to room:", msg.Room)
			pool.mu.RLock()
			for c := range pool.Rooms[msg.Room] {
				if c.UserId == msg.Sender {
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
