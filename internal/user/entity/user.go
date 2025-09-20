package entity

import "database/sql"

// User represents a user in the system.
type User struct {
	ID             sql.NullString `db:"id"`              // Unique ID of the user.
	NickName       sql.NullString `db:"nick_name"`       // Nickname of the user.
	Username       sql.NullString `db:"username"`        // Username of the user.
	HashedPassword sql.NullString `db:"hashed_password"` // Hashed password of the user.
	Enable2FA      sql.NullBool   `db:"enable_2fa"`      // Indicates if 2FA is enabled for the user.
	Email          sql.NullString `db:"email"`           // Email of the user.
	Google         sql.NullString `db:"google"`          // Google of the user.
	Facebook       sql.NullString `db:"facebook"`        // Facebook of the user.
	Discord        sql.NullString `db:"discord"`         // Discord of the user.
	Github         sql.NullString `db:"github"`          // Github of the user.
	Phone          sql.NullString `db:"phone"`           // Phone number of the user.
	OTPSecret      sql.NullString `db:"otp_secret"`      // OTP secret of the user.
}

// TableName returns the table name for the User struct.
func (aa *User) TableName() string {
	return "users"
}
