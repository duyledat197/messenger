package repository

import (
	"context"

	"openmyth/messgener/internal/chat/entity"
)

// ChannelRepository is an interface that defines the contract for channel operations.
type ChannelRepository interface {
	Create(context.Context, *entity.Channel) error
	RetrieveByChannelID(context.Context, int64) (*entity.Channel, error)
	SearchByName(ctx context.Context, name string, offset, limit int64) ([]*entity.Channel, error)
	List(ctx context.Context, offset, limit int64) ([]*entity.Channel, error)
	Delete(context.Context, int64) error
}

type CacheChannelRepository interface {
	CreateByList(ctx context.Context, offset, limit int64, channels []*entity.Channel) error
	List(ctx context.Context, offset, limit int64) ([]*entity.Channel, error)
}
