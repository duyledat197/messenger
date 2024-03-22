package cqlx

import (
	"context"

	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/scylla"
)

// metadata specifies table name and columns it must be in sync with schema.
var messageMetadata = table.Metadata{
	Name:    "message",
	Columns: []string{"user_id", "client_id"},
	PartKey: []string{"user_id"},
	SortKey: []string{},
}

// messageTable allows for simple CRUD operations based on personMetadata.
var messageTable = table.New(messageMetadata)

type messageRepository struct {
	db *scylla.ScyllaClient
}

// NewMessageRepository creates a new instance of the MessageRepository interface.
//
// It returns a pointer to an messageRepository struct that implements the MessageRepository interface.
func NewMessageRepository() repository.MessageRepository {
	return &messageRepository{}
}

// Create creates a new message in the message repository.
//
// ctx: The context.Context object for the request.
// data: The entity.Message object containing the data for the new message.
// Returns an error if there was a problem creating the message.
func (r *messageRepository) Create(ctx context.Context, data *entity.Message) error {
	if err := messageTable.InsertQueryContext(ctx, r.db.Session).BindStruct(data).ExecRelease(); err != nil {
		return err
	}

	return nil
}

// RetrieveByUserID retrieves a list of messages for a given user ID, starting from the specified offset and limited to the specified limit.
func (r *messageRepository) RetrieveByUserID(ctx context.Context, userID string, offset, limit int64) ([]*entity.Message, error) {
	var result []*entity.Message
	if err := messageTable.SelectBuilder().Where(
		qb.Eq("user_id"),
		qb.Eq("bucket"),
		qb.Gt("created_at"),
	).Limit(uint(limit)).QueryContext(ctx, r.db.Session).
		BindStruct(&entity.Message{
			Bucket:    offset,
			UserID:    userID,
			CreatedAt: offset,
		}).SelectRelease(&result); err != nil {
		return nil, err
	}

	return result, nil
}
