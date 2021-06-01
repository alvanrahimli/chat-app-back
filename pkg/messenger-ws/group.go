package messenger_ws

type PrivacyType string

const (
	Private PrivacyType = "private"
	Public 	PrivacyType = "public"
	Hidden 	PrivacyType = "unlisted"
)

type Group struct {
	ID		int
	Name	string
	Privacy	PrivacyType
	Clients []Client
}
