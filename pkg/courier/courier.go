// Package courier represents of courier implementation
package courier

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gojek/courier-go"

	"openmyth/messgener/config"
)

// Client represents a wrapper around a courier.Client with the additional config.Database
// This struct is used to hold the configuration and the courier client.
//
// cfg - holds the configuration for the courier client
// Client - holds the actual courier client
type Client struct {
	*courier.Client                  // Holds the actual courier client
	cfg             *config.Database // Holds the configuration for the courier client
}

// NewClient creates a new Client with the given configuration.
// It takes a pointer to a config.Database struct as a parameter.
// It returns a pointer to a Client struct.
func NewClient(cfg *config.Database) *Client {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatal("courier port is not correct")
	}
	client, err := courier.NewClient(
		courier.WithUsername(cfg.User),
		courier.WithPassword(cfg.Password),
		courier.WithAddress(cfg.Host, uint16(port)),
	)

	if err != nil {
		log.Fatalf("unable to create courier client: %v", err)
	}

	return &Client{
		Client: client,
	}
}

// Connect establishes a connection to the client.
// It starts the client and returns an error if it fails to start.
// Returns nil if the connection is established successfully.
func (c *Client) Connect(_ context.Context) error {
	if err := c.Client.Start(); err != nil {
		return fmt.Errorf("unable to start courier: %w", err)
	}

	return nil
}

// Close stops the client and returns an error if any.
func (c *Client) Close(_ context.Context) error {
	c.Client.Stop()

	return nil
}
