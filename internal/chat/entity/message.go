package entity

type Message struct {
	Bucket    int64    `partkey:"bucket"`
	ChannelID int64    `partkey:"channel_id"`
	MessageID int64    `sortkey:"message_id"`
	UserID    string   `json:"user_id"`
	Content   string   `json:"content"`
	Reaction  []string `json:"reaction"`
}

// TableName returns the name of the table for the Message struct.
//
// This function does not take any parameters.
// It returns a string representing the table name for the Message struct.
func (aa *Message) TableName() string {
	return "messages"
}
