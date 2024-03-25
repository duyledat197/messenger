package repository

import (
	"context"

	"openmyth/messgener/internal/chat/entity"
)

// MessageRepository defines the contract for a repository that handles message operations.
type MessageRepository interface {
	Create(ctx context.Context, data *entity.Message) error
	RetrieveMessages(ctx context.Context, channelID, offset, limit int64) ([]*entity.Message, error)
}
