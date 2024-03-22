package repository

import (
	"context"

	"openmyth/messgener/internal/chat/entity"
)

// OnlineRepository is an interface that defines the contract for retrieving an online status by user ID.
type OnlineRepository interface {
	RetrieveByUserID(ctx context.Context, userID string) (*entity.Online, error)
	Create(ctx context.Context, data *entity.Online) error
	DeleteByUserID(ctx context.Context, userID string) error
}
