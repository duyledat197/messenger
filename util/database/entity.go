package database

// e is a presentation of a entity that must have TableName function inside.
type Entity interface {
	TableName() string
}
