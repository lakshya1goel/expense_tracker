package ws

import (
	"fmt"
	"sync"

	"github.com/lakshya1goel/expense_tracker/database"
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

		case msg := <-pool.Broadcast:
			fmt.Println("Broadcasting message to group:", msg.GroupID)
			fmt.Printf("msg: %+v\n", msg)
			if err := database.Db.Create(&msg).Error; err != nil {
				fmt.Println("DB insert error:", err)
			}
			
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

		// case msg := <-pool.Broadcast:
		// 	fmt.Println("Broadcasting message to group:", msg.GroupID)
		// 	fmt.Printf("msg: %+v\n", msg)

		// 	pool.mu.RLock()
		// 	clients := pool.Groups[msg.GroupID]
		// 	pool.mu.RUnlock()

		// 	switch msg.Type {
		// 		case models.ChatMessage:
		// 			for c := range clients {
		// 				if c.UserId == msg.Sender {
		// 					continue
		// 				}
		// 				if err := c.Conn.WriteJSON(msg); err != nil {
		// 					fmt.Println("Broadcast write error:", err)
		// 				}
		// 			}
		// 			if err := database.Db.Create(&msg).Error; err != nil {
		// 				fmt.Println("DB insert error:", err)
		// 			}
		// 		case models.SplitMessage:
		// 			var split models.Split
		// 			err := json.Unmarshal([]byte(msg.Body), &split)
		// 			if err != nil {
		// 				fmt.Println("Invalid split data:", err)
		// 				break
		// 			}
		// 			split.GroupID = msg.GroupID
		// 			if err := database.Db.Create(&split).Error; err != nil {
		// 				fmt.Println("DB insert error (split):", err)
		// 			}
		// 			for c := range clients {
		// 				if c.UserId == msg.Sender {
		// 					continue
		// 				}
		// 				if err := c.Conn.WriteJSON(msg); err != nil {
		// 					fmt.Println("Broadcast write error:", err)
		// 				}
		// 			}
		// 		default:
		// 			fmt.Println("Unknown message type in broadcast")
		// 	}
		}
	}
}
