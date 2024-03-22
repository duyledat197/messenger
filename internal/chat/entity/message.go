package entity

type Message struct {
	Bucket    int64
	MessageID int64
	UserID    string
	FromID    string
	ToID      string
	Content   string
	Reaction  []string
}
