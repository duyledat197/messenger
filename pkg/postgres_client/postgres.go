// Package ...
package postgres_client

import (
	"context"
	"database/sql"
	"log"
	"openmyth/messgener/config"

	_ "github.com/lib/pq"
)

// PostgresClient is presentation for a custom client of postgres with [database/sql] based.
type PostgresClient struct {
	*sql.DB
	cfg *config.Database
}

// NewPostgresClient creates a new Postgres client with the given connection string.
//
// connString string
// *PostgresClient
func NewPostgresClient(cfg *config.Database) *PostgresClient {
	return &PostgresClient{
		cfg: cfg,
	}
}

// Connect establishes a connection to the PostgreSQL database using the provided connection string.
func (c *PostgresClient) Connect(_ context.Context) error {
	var err error
	c.DB, err = sql.Open("postgres", c.cfg.Address())
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
