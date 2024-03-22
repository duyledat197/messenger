package chat

import "time"

// Message represents a chat message.
type Message struct {
	To        string    // The recipient of the message.
	Content   string    // The actual text of the message.
	Timestamp time.Time // The time the message was sent.
}
