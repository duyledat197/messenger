package scylla

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"

	"openmyth/messgener/config"
)

// ScyllaClient represents a Scylla database session.
// It embeds a `gocqlx.Session` and holds a `gocql.ClusterConfig` config.
type ScyllaClient struct {
	// Session holds the gocqlx session.
	// It embeds a `gocqlx.Session`.
	gocqlx.Session

	// cfg holds the Scylla cluster configuration.
	// It represents the `gocql.ClusterConfig` config.
	cfg *gocql.ClusterConfig
}

// NewScylla creates a new Scylla session.
//
// It takes a string parameter `addr` which represents the address of the Scylla cfg.
// The function returns a pointer to a `gocqlx.Session` object.
func NewScylla(config *config.Database) *ScyllaClient {

	cfg := gocql.NewCluster(config.Address())
	cfg.ConnectTimeout = 2 * time.Second
	cfg.Timeout = time.Second
	cfg.NumConns = 5
	cfg.Logger = log.Default()

	return &ScyllaClient{
		cfg: cfg,
	}
}

// Connect initializes and starts the ScyllaDB.
func (db *ScyllaClient) Connect(_ context.Context) error {
	session, err := gocqlx.WrapSession(db.cfg.CreateSession())
	if err != nil {
		return fmt.Errorf("unable to connect scylladb: %v", err)
	}

	db.Session = session

	return nil
}

// Close stops the ScyllaClient.
func (db *ScyllaClient) Close(_ context.Context) error {
	defer db.Session.Close()

	return nil
}
