package messenger_ws

import (
	"log"
	"time"
)

type Pool struct {
	StartTime	time.Time

	Register	chan *Client
	Unregister	chan *Client
	MakeOnline	chan string
	Clients		map[*Client]bool
	Groups		[]*Group
	Send		chan Message
}

func InitializePool() *Pool {
	return &Pool{
		StartTime: time.Now().UTC(),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		MakeOnline: make(chan string),
		Clients: 	make(map[*Client]bool),
		Groups:     make([]*Group, 10),
		Send:       make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <- pool.Register:
			pool.Clients[client] = true
			log.Printf("New client (%s) registered", client.ID)
			log.Printf("Current pool size: %d", len(pool.Clients))

		case client := <- pool.Unregister:
			log.Printf("Client (%s) unregistered", client.ID)
			delete(pool.Clients, client)

		case id := <- pool.MakeOnline:
			for client := range pool.Clients {
				if client.ID == id {
					pool.Clients[client] = true
				}
			}


		case message := <- pool.Send:
			log.Println(message.Content)
		}
	}
}