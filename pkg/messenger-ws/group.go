package messenger_ws

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type PrivacyType string

const (
	Private PrivacyType = "private"
	Public 	PrivacyType = "public"
	Hidden 	PrivacyType = "unlisted"
)

type Group struct {
	ID			string
	Name		string
	Privacy		PrivacyType
	Clients 	[]*Client
}

type CreateGroupContext struct {
	Name		string `json:"name"`
	Privacy 	PrivacyType `json:"privacy"`
}

func (group *Group) HandleNewMessage(senderClient *Client, message *NewMessageContext) {
	for _, c := range group.Clients {				// Iterate clients of group
		if c.ID != senderClient.ID {				// If it is other members
			response := ClientResponse{
				Status:  Ok,
				Type: NewMessageRes,
				Content: NewMessageResponse{
					Sender:    senderClient.Name,
					Content:   message.Content,
					Timestamp: time.Now().UTC(),
				},
			}

			sendErr := senderClient.Send(response)
			if sendErr != nil {
				log.Printf("Error: %s", sendErr.Error())
				return
			}
		}
	}
}

func (createGroupContext *CreateGroupContext) HandleRequest(client *Client) {
	newGroup := Group{
		ID:      uuid.New().String(),
		Name:    createGroupContext.Name,
		Privacy: createGroupContext.Privacy,
		Clients: []*Client{client},
	}

	client.Pool.Groups = append(client.Pool.Groups, &newGroup)
}