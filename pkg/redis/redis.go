package redis

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	redis "github.com/redis/go-redis/v9"

	"openmyth/messgener/config"
)

type Client struct {
	Client *redis.Client
	cfg    *config.Database
}

func NewClient(cfg *config.Database) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:           net.JoinHostPort(cfg.Host, cfg.Port),
		Password:       cfg.Password, // no password set
		DB:             0,            // use default DB
		MaxActiveConns: int(cfg.MaxConnection),
	})
	return &Client{
		Client: client,
		cfg:    cfg,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	if cmd := c.Client.Ping(ctx); cmd.Err() != nil {
		return fmt.Errorf("unable to connect redis: %w", cmd.Err())
	}

	slog.Info("connect redis successful")
	return nil
}

func (c *Client) Close(_ context.Context) error {
	slog.Info("close redis successful")
	return nil
}
