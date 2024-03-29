package entity

import "database/sql"

// User represents a user in the system.
type User struct {
	ID             sql.NullString `json:"id,omitempty"`              // Unique ID of the user.
	NickName       sql.NullString `json:"nick_name,omitempty"`       // Nickname of the user.
	Username       sql.NullString `json:"username,omitempty"`        // Username of the user.
	HashedPassword sql.NullString `json:"hashed_password,omitempty"` // Hashed password of the user.
	Enable2FA      sql.NullBool   `json:"enable_2fa,omitempty"`      // Indicates if 2FA is enabled for the user.
	Email          sql.NullString `json:"email,omitempty"`           // Email of the user.
	Google         sql.NullString `json:"google,omitempty"`          // Google of the user.
	Facebook       sql.NullString `json:"facebook,omitempty"`        // Facebook of the user.
	Discord        sql.NullString `json:"discord,omitempty"`         // Discord of the user.
	Github         sql.NullString `json:"github,omitempty"`          // Github of the user.
}

// TableName returns the table name for the User struct.
func (aa *User) TableName() string {
	return "users"
}
