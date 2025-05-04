package models

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
	Sender uint `json:"sender"`
	Room uint `json:"room"`
}