package grpc_server

import (
	"context"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

// NewHealthService ...
func NewHealthService() pb.HealthServer {
	return &healthService{}
}

type healthService struct{}

// Check is a function that checks the health of the service.
func (s *healthService) Check(context.Context, *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is a function that description of the Go function.
func (s *healthService) Watch(*pb.HealthCheckRequest, pb.Health_WatchServer) error {
	return nil
}
