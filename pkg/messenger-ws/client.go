package messenger_ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	ID			string
	Name		string
	Conn		*websocket.Conn
	LastPing 	time.Time
	Pool		*Pool
}

func (client *Client) Read() {
	defer func() {
		err := client.Conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
		client.Pool.Unregister <- client
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			log.Printf("Error occurred while reading message from %s", client.ID)
			break
		}

		log.Println(string(p))
	}
}