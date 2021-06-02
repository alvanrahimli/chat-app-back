package messenger_ws

import (
	"log"
	"time"
)

type GetClientsContext struct {
	ClientId	string
}

type Pool struct {
	StartTime	time.Time

	Register	chan *Client
	Unregister	chan *Client
	MakeOnline	chan string
	MakeOffline chan string
	Clients		map[*Client]bool
	Groups		[]*Group
	Send		chan NewMessageContext
}

func InitializePool() *Pool {
	return &Pool{
		StartTime: 		time.Now().UTC(),
		Register:   	make(chan *Client),
		Unregister: 	make(chan *Client),
		MakeOnline: 	make(chan string),
		MakeOffline:	make(chan string),
		Clients: 		make(map[*Client]bool),
		Groups:     	make([]*Group, 0),
		Send:       	make(chan NewMessageContext),
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
		case id := <- pool.MakeOffline:
			for client := range pool.Clients {
				if client.ID == id {
					pool.Clients[client] = false
				}
			}

		case message := <- pool.Send:
			log.Println(message.Content)
		}
	}
}

func (pool *Pool) HandleGetClients(client *Client) {
	var clients []string
	for c, _ := range client.Pool.Clients {
		clients = append(clients, c.ID)
	}

	response := ClientResponse{
		Status:  Ok,
		Type:    ClientsList,
		Content: clients,
	}
	sendErr := client.Send(response)
	if sendErr != nil {
		log.Printf("Error: %s", sendErr.Error())
	}
}