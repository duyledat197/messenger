package redis

import (
	"context"
	"log/slog"
	"openmyth/messgener/config"

	redis "github.com/redis/go-redis/v9"
)

type Client struct {
	Client *redis.Client
	cfg    *config.Database
}

func NewClient(cfg *config.Database) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:       cfg.Host + ":" + cfg.Port,
		Username:   cfg.User,
		Password:   cfg.Password, // no password set
		DB:         0,            // use default DB
		MaxRetries: int(cfg.MaxConnection),
	})
	return &Client{
		Client: client,
		cfg:    cfg,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	if cmd := c.Client.Ping(ctx); cmd.Err() != nil {
		return cmd.Err()
	}

	slog.Info("connect redis successful")
	return nil
}

func (c *Client) Close(_ context.Context) error {
	slog.Info("close redis successful")
	return nil
}
