package messenger_ws

import (
	"encoding/json"
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

		var clientCommand ClientCommand
		unmarshallErr := json.Unmarshal(p, &clientCommand)
		if unmarshallErr != nil {
			log.Printf("Error: %s", unmarshallErr.Error())
			continue
		}


		if clientCommand.Type == NewMessageCmd {
			var newMessage NewMessageContext
			msgUnmarshallErr := json.Unmarshal([]byte(clientCommand.Content), &newMessage)
			if msgUnmarshallErr != nil {
				log.Printf("Error: %s", msgUnmarshallErr.Error())
				continue
			}

			newMessage.HandleRequest(client)
		} else if clientCommand.Type ==  CreateGroup {
			var createGroupContext CreateGroupContext
			msgUnmarshallErr := json.Unmarshal([]byte(clientCommand.Content), &createGroupContext)
			if msgUnmarshallErr != nil {
				log.Printf("Error: %s", msgUnmarshallErr.Error())
				continue
			}

			createGroupContext.HandleRequest(client)
		} else if clientCommand.Type == AddMember {
			var addMemberContext AddMemberContext
			unmarshallErr := json.Unmarshal([]byte(clientCommand.Content), &addMemberContext)
			if unmarshallErr != nil {
				log.Printf("Error: %s", unmarshallErr.Error())
				continue
			}

			addMemberContext.HandleRequest(client)
		}

		log.Println(string(p))
	}
}

func (client *Client) Send(v interface{}) error {
	data, jsonErr := json.Marshal(v)
	if jsonErr != nil {
		return jsonErr
	}

	writeErr := client.Conn.WriteMessage(websocket.TextMessage, data)
	if writeErr != nil {
		return writeErr
	}

	return nil
}