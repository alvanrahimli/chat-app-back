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
	CreatorID	string
	Clients 	[]*Client
}

type CreateGroupContext struct {
	Name		string `json:"name"`
	Privacy 	PrivacyType `json:"privacy"`
}

type AddMemberContext struct {
	GuestId		string `json:"guest_id"`
	GroupId		string `json:"group_id"`
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
		ID:      	uuid.New().String(),
		Name:    	createGroupContext.Name,
		Privacy: 	createGroupContext.Privacy,
		CreatorID: 	client.ID,
		Clients: 	[]*Client{client},
	}

	client.Pool.Groups = append(client.Pool.Groups, &newGroup)
}

func (addMemberContext *AddMemberContext) HandleRequest(client *Client) {
	for _, group := range client.Pool.Groups {
		if group.ID == addMemberContext.GroupId && group.CreatorID == client.ID {
			for guest, _ := range client.Pool.Clients {
				if guest.ID == addMemberContext.GuestId {
					group.Clients = append(group.Clients, guest)
				}
			}
		}
	}
}