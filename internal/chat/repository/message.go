package repository

import (
	"context"

	"openmyth/messgener/internal/chat/entity"
)

type MessageRepository interface {
	Create(ctx context.Context, data *entity.Message) error
	RetrieveByUserID(ctx context.Context, userID string, offset, limit int64) ([]*entity.Message, error)
}
