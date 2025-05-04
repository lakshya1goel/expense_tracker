package ws

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Clients    map[*Client]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("total connection pool:- ", len(pool.Clients))
			for k := range pool.Clients {
				fmt.Println(k)
				k.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined"})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("total connection pool:- ", len(pool.Clients))
			for k := range pool.Clients {
				fmt.Println(k)
				k.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected"})
			}
		case msg := <-pool.Broadcast:
			fmt.Println("broadcasting a message")
			for k := range pool.Clients {
				if err := k.Conn.WriteJSON(msg); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
