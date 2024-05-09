package repository

import (
	"context"

	"openmyth/messgener/internal/chat/entity"
)

// ChannelRepository is an interface that defines the contract for channel operations.
type ChannelRepository interface {
	Create(context.Context, *entity.Channel) error
	RetrieveByChannelID(context.Context, int64) (*entity.Channel, error)
	SearchByName(context.Context, string) ([]*entity.Channel, error)
	Delete(context.Context, int64) error
}
