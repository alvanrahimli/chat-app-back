package messenger_ws

import (
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

type GroupDto struct {
	ID			string
	Name		string
	Privacy		PrivacyType
	CreatorID	string
}

type CreateGroupContext struct {
	Name		string `json:"name"`
	Privacy 	PrivacyType `json:"privacy"`
}

type AddMemberContext struct {
	GuestId		string `json:"guest_id"`
	GroupId		string `json:"group_id"`
}

type GetGroupsContext struct {
	ClientId	string `json:"client_id"`
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
