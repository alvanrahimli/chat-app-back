package messenger_ws

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

		log.Printf("MESSAGE: %s", string(p))

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

			client.HandleNewMessage(newMessage)
		} else if clientCommand.Type ==  CreateGroup {
			var createGroupContext CreateGroupContext
			msgUnmarshallErr := json.Unmarshal([]byte(clientCommand.Content), &createGroupContext)
			if msgUnmarshallErr != nil {
				log.Printf("Error: %s", msgUnmarshallErr.Error())
				continue
			}

			client.HandleCreateGroup(createGroupContext)
		} else if clientCommand.Type == AddMember {
			var addMemberContext AddMemberContext
			unmarshallErr := json.Unmarshal([]byte(clientCommand.Content), &addMemberContext)
			if unmarshallErr != nil {
				log.Printf("Error: %s", unmarshallErr.Error())
				continue
			}

			client.HandleAddMember(addMemberContext)
		} else if clientCommand.Type == GetClients {
			var getClientsContext GetClientsContext
			unmarshallErr := json.Unmarshal([]byte(clientCommand.Content), &getClientsContext)
			if unmarshallErr != nil {
				log.Printf("Error: %s", unmarshallErr.Error())
				continue
			}

			client.HandleGetClients(getClientsContext)
		}

		log.Println(string(p))
	}
}

func (client *Client) Send(v ClientResponse) error {
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

func (client *Client) HandleCreateGroup(createGroupContext CreateGroupContext) {
	newGroup := Group{
		ID:      	uuid.New().String(),
		Name:    	createGroupContext.Name,
		Privacy: 	createGroupContext.Privacy,
		CreatorID: 	client.ID,
		Clients: 	[]*Client{client},
	}

	client.Pool.Groups = append(client.Pool.Groups, &newGroup)
	sendErr := client.Send(ClientResponse{
		Status:  Ok,
		Type:    GroupCreated,
		Content: newGroup.ID,
	})
	if sendErr != nil {
		log.Printf("Error: %s", sendErr.Error())
	}
}

func (client *Client) HandleAddMember(addMemberContext AddMemberContext) {
	for _, group := range client.Pool.Groups {
		if group.ID == addMemberContext.GroupId && group.CreatorID == client.ID {
			for guest, _ := range client.Pool.Clients {
				if guest.ID == addMemberContext.GuestId {
					group.Clients = append(group.Clients, guest)
					response := fmt.Sprintf("%s:%s", guest.ID, guest.Name)
					group.Broadcast(ClientAdded, response)
					log.Printf("Guest-Client (%s) added to Group (%s)", guest.ID, group.ID)
					return
				}
			}
		}
	}
}

func (client *Client) HandleGetClients(context GetClientsContext) {
	var clients = make(map[string]string)

	if context.GroupId != "" {
		for _, g := range client.Pool.Groups {
			if g.ID == context.GroupId {
				for _, c := range g.Clients {
					clients[c.ID] = c.Name
				}
			}
		}
	} else {
		for c, _ := range client.Pool.Clients {
			clients[c.ID] = c.Name
		}
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