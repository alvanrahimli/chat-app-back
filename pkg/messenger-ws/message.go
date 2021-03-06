package messenger_ws

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	//Ping		MessageType = "ping"
	// Add more message types
)

type NewMessageContext struct {
	GroupID		string `json:"group_id"`
	Header		string `json:"header"`
	Content		string `json:"content"`
}

type CommandType string
const (
	Registration  	CommandType = "register"
	BecomeOnline  	CommandType = "become_alive"
	NewMessageCmd 	CommandType = "new_message"
	GetClients		CommandType = "get_clients"
	GetGroups		CommandType = "get_groups"
	CreateGroup		CommandType = "create_group"
	AddMember		CommandType = "add_member"
)

type ResponseType string
const (
	Registered		ResponseType = "registered"
	NewMessageRes	ResponseType = "new_message"
	GroupCreated	ResponseType = "group_created"
	ClientsList		ResponseType = "clients_list"
	GroupsList		ResponseType = "groups_list"
	ClientAdded		ResponseType = "client_added"
)

type ResponseStatus string
const (
	Ok		ResponseStatus = "ok"
)

type NewMessageResponse struct {
	ID				string
	Sender		string
	Content		string
	Timestamp	time.Time
}

type ClientCommand struct {
	Type	CommandType `json:"type"`
	Content string `json:"content"`
}

type ClientResponse struct {
	Status	ResponseStatus 	`json:"status"`
	Type	ResponseType 	`json:"type"`
	Content interface{} 	`json:"content"`
}

func (client *Client) HandleNewMessage(newMessageContext NewMessageContext) {
	for _, group := range client.Pool.Groups { 				// Iterate groups
		if group.ID == newMessageContext.GroupID { 			// Find group
			group.Broadcast(NewMessageRes, NewMessageResponse {
				ID: uuid.NewString(),
				Sender:    client.Name,
				Content:   newMessageContext.Content,
				Timestamp: time.Now().UTC(),
			})
		}
	}
}
