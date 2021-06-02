package messenger_ws

import (
	"fmt"
	"github.com/google/uuid"
	"log"
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

func (group *Group) Broadcast(responseType ResponseType, content interface{}) {
	for _, c := range group.Clients {				// Iterate clients of group
		response := ClientResponse{
			Status:  Ok,
			Type: responseType,
			Content: content,
		}

		sendErr := c.Send(response)
		if sendErr != nil {
			log.Printf("Error: %s", sendErr.Error())
			return
		}
		log.Printf("Message sent to ClientID: %s", c.ID)
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
	sendErr := client.Send(ClientResponse{
		Status:  Ok,
		Type:    GroupCreated,
		Content: newGroup.ID,
	})
	if sendErr != nil {
		log.Printf("Error: %s", sendErr.Error())
	}
}

func (addMemberContext *AddMemberContext) HandleRequest(client *Client) {
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