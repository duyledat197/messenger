// Package ...
package postgres_client

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// PostgresClient is presentation for a custom client of postgres with [database/sql] based.
type PostgresClient struct {
	*sql.DB
	addr string
}

// NewPostgresClient creates a new Postgres client with the given connection string.
//
// connString string
// *PostgresClient
func NewPostgresClient(addr string) *PostgresClient {
	return &PostgresClient{
		addr: addr,
	}
}

// Connect establishes a connection to the PostgreSQL database using the provided connection string.
func (c *PostgresClient) Connect(_ context.Context) error {
	var err error
	c.DB, err = sql.Open("postgres", c.addr)
	if err != nil {
		return err
	}

	if err := c.DB.Ping(); err != nil {
		return err
	}

	log.Println("connect postgres successful")

	return nil
}

// Close closes the Postgres client.
func (c *PostgresClient) Close(_ context.Context) error {
	defer c.DB.Close()

	return nil
}
