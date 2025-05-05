package models

type Message struct {
	Type int    `json:"type"` // 0: join_room, 1: chat message
	Body string `json:"body"`
	Sender uint `json:"sender"`
	Room string `json:"room"`
}

const (
    JoinRoom    = 0
    LeaveRoom   = 1
    ChatMessage = 2
)