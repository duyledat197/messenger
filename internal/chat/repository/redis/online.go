package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/redis"
)

const (
	OnlineKeyPrefix = "online:"
)

type onlineRepository struct {
	redisClient *redis.Client
}

func NewOnlineRepository(redisClient *redis.Client) repository.OnlineRepository {
	return &onlineRepository{
		redisClient: redisClient,
	}
}

// RetrieveByUserID retrieves an online repository by user ID.
func (r *onlineRepository) RetrieveByUserID(ctx context.Context, userID string) (*entity.Online, error) {
	cmd := r.redisClient.Client.Get(ctx, userID)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	var online entity.Online
	if err := json.Unmarshal([]byte(cmd.Val()), &online); err != nil {
		return nil, err
	}

	return &online, nil
}

// Create creates a new online repository by user ID.
func (r *onlineRepository) Create(ctx context.Context, data *entity.Online) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	cmd := r.redisClient.Client.Set(ctx,
		fmt.Sprintf("%s:%s", OnlineKeyPrefix, data.UserID),
		string(b),
		0,
	)

	return cmd.Err()

}

// DeleteByUserID deletes an online entry from the repository by user ID.

func (r *onlineRepository) DeleteByUserID(ctx context.Context, userID string) error {
	cmd := r.redisClient.Client.Del(ctx, fmt.Sprintf("%s:%s", OnlineKeyPrefix, userID))

	return cmd.Err()
}
