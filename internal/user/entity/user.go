package entity

// User represents a user in the system.
type User struct {
	ID             string `json:"id,omitempty"`              // Unique ID of the user.
	NickName       string `json:"nick_name,omitempty"`       // Nickname of the user.
	Username       string `json:"username,omitempty"`        // Username of the user.
	HashedPassword string `json:"hashed_password,omitempty"` // Hashed password of the user.
	Enable2FA      bool   `json:"enable_2fa,omitempty"`      // Indicates if 2FA is enabled for the user.
	Email          string `json:"email,omitempty"`           // Email of the user.
	Gmail          string `json:"gmail,omitempty"`           // Gmail of the user.
	Facebook       string `json:"facebook,omitempty"`        // Facebook of the user.
	Discord        string `json:"discord,omitempty"`         // Discord of the user.
	Github         string `json:"github,omitempty"`          // Github of the user.
}
