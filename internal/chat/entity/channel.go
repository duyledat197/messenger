package entity

type Channel struct {
	ChannelID   int64  `json:"channel_id,omitempty"`
	Name        string `json:"name,omitempty"`
	MaxUser     uint32 `json:"max_user,omitempty"`
	Description string `json:"description,omitempty"`
}

// TableName returns the table name for the Channel struct.
//
// This function does not take any parameters.
// It returns a string representing the table name for the Channel struct.
func (aa *Channel) TableName() string {
	return "channels"
}
