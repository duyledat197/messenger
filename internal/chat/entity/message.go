package entity

type Message struct {
	Bucket    int64    `partkey:"bucket"`
	ChannelID int64    `partkey:"channel_id"`
	MessageID int64    `sortkey:"message_id"`
	UserID    string   `json:"user_id"`
	Content   string   `json:"content"`
	Reaction  []string `json:"reaction"`
}

func (aa *Message) TableName() string {
	return "messages"
}
