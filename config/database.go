package config

import "fmt"

// Database represents the configuration details for a database connection.
type Database struct {
	Schema        string `mapstructure:"schema"`
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
		"%s://%s:%s@%s:%s/%s",
		e.Schema,
		e.User,
		e.Password,
		e.Host,
		e.Port,
		e.Database,
	)
}
