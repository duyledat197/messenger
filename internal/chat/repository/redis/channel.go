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
	ChannelKeyPrefix = "channel:"
	ChannelListTTL   = 60 * 60
)

type cacheChannelRepository struct {
	redisClient *redis.Client
}

func NewCacheChannelRepository(redisClient *redis.Client) repository.CacheChannelRepository {
	return &cacheChannelRepository{
		redisClient: redisClient,
	}
}

func (c *cacheChannelRepository) List(ctx context.Context, offset, limit int64) ([]*entity.Channel, error) {
	cmd := c.redisClient.Client.Get(ctx, fmt.Sprintf("%s:%d:%d", ChannelKeyPrefix, offset, limit))
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	var result []*entity.Channel
	if err := json.Unmarshal([]byte(cmd.Val()), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *cacheChannelRepository) CreateByList(ctx context.Context, offset, limit int64, channels []*entity.Channel) error {
	data, err := json.Marshal(channels)
	if err != nil {
		return err
	}
	cmd := c.redisClient.Client.Set(ctx, fmt.Sprintf("%s:%d:%d", ChannelKeyPrefix, offset, limit), data, ChannelListTTL)

	return cmd.Err()
}
