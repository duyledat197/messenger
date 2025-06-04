package config

import "fmt"

// Database represents the configuration details for a database connection.
type Database struct {
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	Database      string `mapstructure:"database"`
	MaxConnection int32  `mapstructure:"max_connection"`
}

// Address returns the formatted string for the database connection address.
func (e *Database) Address() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		e.Host,
		e.Port,
		e.User,
		e.Password,
		e.Database,
	)
}
