package cqlx

import (
	"context"

	"github.com/scylladb/gocqlx/v2/table"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/scylla"
	"openmyth/messgener/util/cqlx_util"
	"openmyth/messgener/util/database"
)

type onlineRepository struct {
	db  *scylla.ScyllaClient
	tbl *table.Table
}

// NewOnlineRepository creates a new instance of the OnlineRepository interface.
// It returns a pointer to an onlineRepository struct that implements the OnlineRepository interface.
func NewOnlineRepository(db *scylla.ScyllaClient) repository.OnlineRepository {
	e := &entity.Online{}
	partKeys, sortKeys := cqlx_util.Fields(e)
	fields, _ := database.FieldMap(e)

	metadata := table.Metadata{
		Name:    e.TableName(),
		Columns: fields,
		PartKey: partKeys,
		SortKey: sortKeys,
	}

	return &onlineRepository{
		db:  db,
		tbl: table.New(metadata),
	}
}

// RetrieveByUserID retrieves an online repository by user ID.
func (r *onlineRepository) RetrieveByUserID(ctx context.Context, userID string) (*entity.Online, error) {
	var result entity.Online
	if err := r.tbl.
		GetQueryContext(ctx, r.db.Session).
		BindStruct(&entity.Online{UserID: userID}).
		GetRelease(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create creates a new record in the onlineRepository.
func (r *onlineRepository) Create(ctx context.Context, data *entity.Online) error {
	if err := r.tbl.
		InsertQueryContext(ctx, r.db.Session).
		BindStruct(data).
		ExecRelease(); err != nil {
		return err
	}

	return nil
}

// DeleteByUserID deletes an entry from the online repository based on the userID.
func (r *onlineRepository) DeleteByUserID(ctx context.Context, userID string) error {
	if err := r.tbl.
		DeleteQueryContext(ctx, r.db.Session).
		BindStruct(&entity.Online{UserID: userID}).
		ExecRelease(); err != nil {
		return err
	}

	return nil
}
