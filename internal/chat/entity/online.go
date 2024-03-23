package entity

type Online struct {
	UserID   string
	ClientID string
}

// TableName returns the name of the table for the Online struct.
func (aa *Online) TableName() string {
	return "onlines"
}
