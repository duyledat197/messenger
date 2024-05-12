package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"

	"openmyth/messgener/config"
)

// GrpcServer represents a gRPC server instance.
//
// endpoint *config.Endpoint: the endpoint configuration.
// Server *grpc.Server: the gRPC server.
type GrpcServer struct {
	endpoint     *config.Endpoint // The endpoint configuration.
	*grpc.Server                  // The gRPC server.
}

// NewGrpcServer creates a new instance of GrpcServer.
//
// It takes an endpoint of type *config.Endpoint as a parameter and returns a pointer to GrpcServer.
func NewGrpcServer(endpoint *config.Endpoint) *GrpcServer {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// Start starts the gRPC server.
			//
			// It listens on the specified port and starts serving incoming requests.
			// The function takes a context.Context as a parameter and returns an error if there was a problem starting the server.
			grpc_validator.StreamServerInterceptor(),
		)),
	)
	return &GrpcServer{
		endpoint: endpoint,
		Server:   srv,
	}
}

// Start starts the gRPC server.
//
// It listens on the specified port and starts serving incoming requests.
// The function takes a context.Context as a parameter and returns an error if there was a problem starting the server.
func (s *GrpcServer) Start(_ context.Context) error {
	// Start starts the gRPC server.
	//
	// It listens on the specified port and starts serving incoming requests.
	// The function takes a context.Context as a parameter and returns an error if there was a problem starting the server.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.endpoint.Port))
	if err != nil {
		return err
	}
	log.Printf("Server listening in port: %s\n", s.endpoint.Port)
	if err := s.Server.Serve(lis); err != nil {
		return err
	}

	return nil
}

// Stop stops the GrpcServer gracefully.
func (s *GrpcServer) Stop(_ context.Context) error {
	s.Server.GracefulStop()

	return nil
}
