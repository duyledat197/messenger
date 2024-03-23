package entity

type Message struct {
	Bucket    int64
	ChannelID int64
	MessageID int64
	UserID    string
	Content   string
	Reaction  []string
}

func (aa *Message) TableName() string {
	return "messages"
}
