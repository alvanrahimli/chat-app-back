package messenger_ws

import (
	"time"
)

type MessageType string

const (
	Ping		MessageType = "ping"
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
	CreateGroup		CommandType = "create_group"
	AddMember		CommandType = "add_member"
)

type ResponseType string
const (
	NewMessageRes	ResponseType = "new_message"
	GroupCreated	ResponseType = "group_created"
)

type ResponseStatus string
const (
	Ok		ResponseStatus = "ok"
)

type NewMessageResponse struct {
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

func (newMessageContext *NewMessageContext) HandleRequest(client *Client) {
	for _, group := range client.Pool.Groups { 				// Iterate groups
		if group.ID == newMessageContext.GroupID { // Find group
			group.HandleNewMessage(client, newMessageContext)
		}
	}
}
