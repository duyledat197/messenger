package grpc_client

import (
	"context"
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"openmyth/messgener/config"
)

// GrpcClient represents a gRPC client connection.
//
// *grpc.ClientConn: the gRPC client connection.
// cfg: the configuration for the endpoint.
type GrpcClient struct {
	*grpc.ClientConn                  // The gRPC client connection.
	cfg              *config.Endpoint // The configuration for the endpoint.
}

// NewGrpcClient creates a new GrpcClient with the given config endpoint.
//
// cfg: the config endpoint for the GrpcClient.
// *GrpcClient: returns a pointer to the newly created GrpcClient.
func NewGrpcClient(cfg *config.Endpoint) *GrpcClient {
	return &GrpcClient{
		cfg: cfg,
	}
}

// Connect establishes a connection to the gRPC server.
//
// ctx: the context to use for the connection.
// error: returns an error if the connection fails.
func (c *GrpcClient) Connect(ctx context.Context) error {
	optsRetry := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(50 * time.Millisecond)),
		grpc_retry.WithPerRetryTimeout(3 * time.Second),
	}

	conn, err := grpc.NewClient(
		c.cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_retry.StreamClientInterceptor(optsRetry...),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(optsRetry...),
		)),
	)

	if err != nil {
		return err
	}

	log.Println("connect grpc successful, address:", c.cfg.Host+":"+c.cfg.Port)
	c.ClientConn = conn

	return nil
}

// Close closes the gRPC client connection.
//
// It takes a context.Context as a parameter and returns an error if there was a problem closing the connection.
func (c *GrpcClient) Close(ctx context.Context) error {
	defer c.ClientConn.Close()

	log.Println("close grpc successful, address: ", c.cfg.Host+":"+c.cfg.Port)
	return nil
}
