package scylla

import (
	"context"

	"github.com/scylladb/gocqlx/v2/migrate"

	"openmyth/messgener/migration/cql"
)

// Migrate migrates the ScyllaClient to the latest version.
func (db *ScyllaClient) Migrate(_ context.Context) error {

	if err := migrate.FromFS(context.Background(), db.Session, cql.Files); err != nil {
		return err
	}

	// Second run skips the processed files
	if err := migrate.FromFS(context.Background(), db.Session, cql.Files); err != nil {
		return err
	}

	return nil
}
