package entity

import "time"

type Online struct {
	UserID       string
	ClientID     string
	LastActiveAt time.Time
}

// TableName returns the table name for the Online struct.
//
// This function does not take any parameters.
// It returns a string representing the table name for the Online struct.
func (aa *Online) TableName() string {
	return "onlines"
}
