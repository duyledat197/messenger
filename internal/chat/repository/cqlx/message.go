package cqlx

import (
	"context"

	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/scylla"
	"openmyth/messgener/util/cqlx_util"
	"openmyth/messgener/util/database"
	"openmyth/messgener/util/snowflake"
)

type messageRepository struct {
	db *scylla.ScyllaClient

	tbl *table.Table
}

// NewMessageRepository creates a new instance of the MessageRepository interface.
//
// It returns a pointer to an messageRepository struct that implements the MessageRepository interface.
func NewMessageRepository(db *scylla.ScyllaClient) repository.MessageRepository {
	e := &entity.Message{}
	partKeys, sortKeys := cqlx_util.Fields(e)
	fields, _ := database.FieldMap(e)
	metadata := table.Metadata{
		Name:    e.TableName(),
		PartKey: partKeys,
		SortKey: sortKeys,
		Columns: fields,
	}
	return &messageRepository{
		db,
		table.New(metadata),
	}
}

// Create creates a new message in the message repository.
//
// ctx: The context.Context object for the request.
// data: The entity.Message object containing the data for the new message.
// Returns an error if there was a problem creating the message.
func (r *messageRepository) Create(ctx context.Context, data *entity.Message) error {
	if err := r.tbl.InsertQueryContext(ctx, r.db.Session).BindStruct(data).ExecRelease(); err != nil {
		return err
	}

	return nil
}

// RetrieveByUserID retrieves a list of messages for a given user ID, starting from the specified offset and limited to the specified limit.
func (r *messageRepository) RetrieveMessages(ctx context.Context, channelID, offset, limit int64) ([]*entity.Message, error) {
	buckets := snowflake.MakeBuckets(channelID, offset)
	var result []*entity.Message
	if err := r.tbl.
		SelectBuilder().
		Where(
			qb.Eq("channel_id"),
			qb.In("buckets"),
			qb.Gt("message_id"),
		).
		LimitPerPartition(uint(limit)).
		Limit(uint(limit)).
		QueryContext(ctx, r.db.Session).
		BindMap(map[string]any{
			"channel_id": channelID,
			"buckets":    buckets,
			"message_id": offset,
		}).
		SelectRelease(&result); err != nil {
		return nil, err
	}

	return result, nil
}
