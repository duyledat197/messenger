package websocket

// Impl is an interface that represents a websocket implementation.
//
// It has two methods, Write and Read, both of which take a generic type T as their argument and return an error. The Write method writes a message of type T to the websocket connection, and the Read method reads a message of type T from the websocket connection.
type Impl[T any] interface {
	Execute(client *Client[T], data T) error
	Register(clientID, userID string) error
	UnRegister(userID string) error
}
