package entity

type Online struct {
	UserID   string `partkey:"user_id"`
	ClientID string
}

// TableName returns the table name for the Online struct.
//
// This function does not take any parameters.
// It returns a string representing the table name for the Online struct.
func (aa *Online) TableName() string {
	return "onlines"
}
