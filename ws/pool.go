package ws

import (
	"fmt"
	"sync"

	"github.com/lakshya1goel/expense_tracker/models"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Groups      map[uint]map[*Client]bool
	Broadcast  chan models.Message
	mu         sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Groups:      make(map[uint]map[*Client]bool),
		Broadcast:  make(chan models.Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			if client.GroupID == 0 {
				fmt.Println("Error: Client has not joined a valid group. Skipping registration.")
				continue
			}

			pool.mu.Lock()
			if _, ok := pool.Groups[client.GroupID]; !ok {
				pool.Groups[client.GroupID] = make(map[*Client]bool)
			}
			pool.Groups[client.GroupID][client] = true
			pool.mu.Unlock()

			fmt.Printf("New user joined group: %s, total users: %d\n", client.GroupID, len(pool.Groups[client.GroupID]))

			for c := range pool.Groups[client.GroupID] {
				if c != client {
					err := c.Conn.WriteJSON(models.Message{
						Type:   models.JoinGroup,
						Body:   "New User Joined",
						Sender: 0,
						GroupID:   client.GroupID,
					})
					if err != nil {
						fmt.Println("Write error:", err)
					}
				}
			}
		case client := <-pool.Unregister:
			if client.GroupID == 0 {
				fmt.Println("Error: Client with invalid GroupID tried to unregister.")
				continue
			}

			pool.mu.Lock()
			if _, ok := pool.Groups[client.GroupID]; ok {
				delete(pool.Groups[client.GroupID], client)
				if len(pool.Groups[client.GroupID]) == 0 {
					delete(pool.Groups, client.GroupID)
				}
			}
			pool.mu.Unlock()

			for c := range pool.Groups[client.GroupID] {
				err := c.Conn.WriteJSON(models.Message{
					Type:   models.LeaveGroup,
					Body:   "User Disconnected",
					Sender: 0,
					GroupID: client.GroupID,
				})
				if err != nil {
					fmt.Println("Write error:", err)
				}
			}
		case msg := <-pool.Broadcast:
			fmt.Println("Broadcasting message to group:", msg.GroupID)
			pool.mu.RLock()
			for c := range pool.Groups[msg.GroupID] {
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
