// Package websocket ...
package websocket

import (
	"context"
	"log/slog"
)

type SendInfo[T any] struct {
	ToClientID string
	Data       T
}

// Engine maintains the set of active clients and their message/registration channels.
// It's responsible for broadcasting messages to the clients.
type Engine[T any] struct {
	// Registered clients.
	clients map[string]*Client[T] // currently connected clients.

	// Inbound messages from the clients.
	broadcast chan T // broadcast messages to the clients.

	// Register requests from the clients.
	register chan *Client[T] // register a client to the hub.

	// Unregister requests from clients.
	unregister chan *Client[T] // unregister a client from the hub.

	Send chan SendInfo[T]
}

// NewEngine creates a new Engine and returns a pointer to it.
func NewEngine[T any]() *Engine[T] {
	return &Engine[T]{
		broadcast:  make(chan T),
		register:   make(chan *Client[T]),
		unregister: make(chan *Client[T]),
		clients:    make(map[string]*Client[T]),
	}
}

// Start starts the Engine and handles incoming client connections, disconnections, and broadcast messages.
func (h *Engine[T]) Start(_ context.Context) error {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
		case sendInfo := <-h.Send:
			if cli, ok := h.clients[sendInfo.ToClientID]; ok {
				cli.Send <- sendInfo.Data
			} else {
				slog.Error("client ID not found")
			}
		case message := <-h.broadcast:
			for clientID := range h.clients {
				client := h.clients[clientID]
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
		}
	}
}

// Stop description of the Go function.
func (h *Engine[T]) Stop(_ context.Context) error {
	close(h.register)
	close(h.unregister)
	close(h.broadcast)
	close(h.Send)

	for _, client := range h.clients {
		close(client.Send)
		delete(h.clients, client.ID)
	}

	return nil
}
