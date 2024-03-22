package cqlx

import (
	"context"

	"github.com/scylladb/gocqlx/v2/table"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/scylla"
)

// metadata specifies table name and columns it must be in sync with schema.
var onlineMetadata = table.Metadata{
	Name:    "online",
	Columns: []string{"user_id", "client_id"},
	PartKey: []string{"user_id"},
	SortKey: []string{},
}

// onlineTable allows for simple CRUD operations based on personMetadata.
var onlineTable = table.New(onlineMetadata)

type onlineRepository struct {
	db *scylla.ScyllaClient
}

// NewOnlineRepository creates a new instance of the OnlineRepository interface.
//
// It returns a pointer to an onlineRepository struct that implements the OnlineRepository interface.
func NewOnlineRepository() repository.OnlineRepository {
	return &onlineRepository{}
}

// RetrieveByUserID retrieves an online repository by user ID.
func (r *onlineRepository) RetrieveByUserID(ctx context.Context, userID string) (*entity.Online, error) {
	var result entity.Online
	if err := onlineTable.
		GetQueryContext(ctx, r.db.Session).
		BindStruct(&entity.Online{UserID: userID}).
		GetRelease(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create creates a new record in the onlineRepository.
//
// ctx: context for the operation.
// data: the entity.Online data to be inserted.
// error: returns an error if the operation fails.
func (r *onlineRepository) Create(ctx context.Context, data *entity.Online) error {
	if err := onlineTable.InsertQueryContext(ctx, r.db.Session).BindStruct(data).ExecRelease(); err != nil {
		return err
	}

	return nil
}

// DeleteByUserID deletes an entry from the online repository based on the userID.
//
// ctx is the context for the operation.
// userID is the identifier for the user.
// Returns an error if any.
func (r *onlineRepository) DeleteByUserID(ctx context.Context, userID string) error {
	if err := onlineTable.DeleteQueryContext(ctx, r.db.Session).BindStruct(&entity.Online{UserID: userID}).ExecRelease(); err != nil {
		return err
	}

	return nil
}
