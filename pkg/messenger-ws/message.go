package messenger_ws

type MessageType string

const (
	NewMessage	MessageType = "new_message"
	Ping		MessageType = "ping"
	// Add more message types
)

type Message struct {
	Header		string `json:"header"`
	Content		string `json:"content"`
}
